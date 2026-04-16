import * as THREE from 'three';
import { Pattern, Part, Round } from '../types';

/**
 * Generates a 3D mesh from crochet pattern instructions
 * Creates actual crocheted topology with rounds and stitches
 */
export class CrochetMeshGenerator {
  /**
   * Generate a mesh from a pattern
   */
  static generateFromPattern(pattern: Pattern): THREE.BufferGeometry | null {
    if (!pattern.parts || pattern.parts.length === 0) {
      return null;
    }

    // For now, handle the first part (usually the body)
    const part = pattern.parts[0];

    // Determine part type and generate appropriate mesh
    if (part.type === 'sphere') {
      return this.generateSphereMesh(part);
    }

    // Default to sphere if unknown type
    return this.generateSphereMesh(part);
  }

  /**
   * Generate a spherical crocheted mesh from rounds
   */
  private static generateSphereMesh(part: Part): THREE.BufferGeometry {
    const rounds = part.rounds.filter(r => r.stitch_count > 0); // Exclude finish rounds

    if (rounds.length === 0) {
      // Fallback to simple sphere
      return new THREE.SphereGeometry(1, 12, 12);
    }

    const vertices: number[] = [];
    const indices: number[] = [];
    const normals: number[] = [];
    const uvs: number[] = [];

    // Calculate total height based on number of rounds
    const totalRounds = rounds.length;
    const roundHeight = 1.0 / totalRounds; // Unit height per round

    // Add top center point (magic ring center)
    vertices.push(0, 0.5, 0);
    normals.push(0, 1, 0);
    uvs.push(0.5, 1);

    let vertexIndex = 1;

    // Generate vertices for each round
    for (let roundIdx = 0; roundIdx < rounds.length; roundIdx++) {
      const round = rounds[roundIdx];
      const stitchCount = round.stitch_count;

      if (stitchCount === 0) continue;

      // Calculate Y position (height) for this round
      // Start from top (0.5) and work down to bottom (-0.5)
      const yPos = 0.5 - (roundIdx + 1) * roundHeight;

      // Calculate radius at this round
      // Sphere bulges out in middle, smaller at top/bottom
      const normalizedHeight = (roundIdx + 1) / totalRounds;
      let radius: number;

      if (normalizedHeight < 0.5) {
        // Increasing phase (top to equator)
        radius = Math.sin(normalizedHeight * Math.PI);
      } else {
        // Decreasing phase (equator to bottom)
        radius = Math.sin((1 - normalizedHeight) * Math.PI);
      }

      // Scale radius based on stitch count (more stitches = wider)
      const maxStitches = Math.max(...rounds.map(r => r.stitch_count));
      radius *= (stitchCount / maxStitches) * 0.8;

      // Add slight variation for hand-crocheted look
      const bumpiness = 0.05;

      // Create vertices around this round
      for (let i = 0; i < stitchCount; i++) {
        const angle = (i / stitchCount) * Math.PI * 2;

        // Add slight irregularity to simulate hand-crocheted texture
        const radiusVariation = radius + Math.sin(i * 2.5) * bumpiness;

        const x = Math.cos(angle) * radiusVariation;
        const z = Math.sin(angle) * radiusVariation;
        const y = yPos + Math.sin(i * 3) * bumpiness * 0.5; // Slight height variation

        vertices.push(x, y, z);

        // Calculate normal (pointing outward)
        const normal = new THREE.Vector3(x, y, z).normalize();
        normals.push(normal.x, normal.y, normal.z);

        // UV coordinates
        uvs.push(i / stitchCount, 1 - normalizedHeight);
      }

      // Create faces connecting to previous round
      if (roundIdx === 0) {
        // Connect first round to center point
        const currentRoundStart = 1;
        for (let i = 0; i < stitchCount; i++) {
          const next = (i + 1) % stitchCount;
          indices.push(
            0, // Center point
            currentRoundStart + next,
            currentRoundStart + i
          );
        }
      } else {
        // Connect this round to previous round
        const prevRound = rounds[roundIdx - 1];
        const prevStitchCount = prevRound.stitch_count;
        const prevRoundStart = vertexIndex - prevStitchCount;
        const currentRoundStart = vertexIndex;

        // Create faces between rounds
        // Handle different stitch counts (increases/decreases)
        const ratio = stitchCount / prevStitchCount;

        if (ratio >= 1) {
          // Increasing or same (each prev stitch connects to one or more current stitches)
          for (let i = 0; i < prevStitchCount; i++) {
            const prevIdx = prevRoundStart + i;
            const prevNext = prevRoundStart + ((i + 1) % prevStitchCount);

            const currStart = Math.floor(i * ratio);
            const currEnd = Math.floor((i + 1) * ratio);

            for (let j = currStart; j < currEnd; j++) {
              const currIdx = currentRoundStart + (j % stitchCount);
              const currNext = currentRoundStart + ((j + 1) % stitchCount);

              // Create quad as two triangles
              indices.push(prevIdx, currNext, currIdx);
              indices.push(prevIdx, prevNext, currNext);
            }
          }
        } else {
          // Decreasing (multiple prev stitches to one current stitch)
          for (let i = 0; i < stitchCount; i++) {
            const currIdx = currentRoundStart + i;
            const currNext = currentRoundStart + ((i + 1) % stitchCount);

            const prevStart = Math.floor(i / ratio);
            const prevEnd = Math.floor((i + 1) / ratio);

            for (let j = prevStart; j <= prevEnd && j < prevStitchCount; j++) {
              const prevIdx = prevRoundStart + (j % prevStitchCount);
              const prevNext = prevRoundStart + ((j + 1) % prevStitchCount);

              indices.push(prevIdx, currNext, currIdx);
              if (j < prevEnd) {
                indices.push(prevIdx, prevNext, currNext);
              }
            }
          }
        }
      }

      vertexIndex += stitchCount;
    }

    // Add bottom center point if last round has stitches
    const lastRound = rounds[rounds.length - 1];
    if (lastRound.stitch_count > 0) {
      const lastRoundStart = vertexIndex - lastRound.stitch_count;

      vertices.push(0, -0.5, 0);
      normals.push(0, -1, 0);
      uvs.push(0.5, 0);

      const bottomIdx = vertexIndex;

      // Connect last round to bottom center
      for (let i = 0; i < lastRound.stitch_count; i++) {
        const next = (i + 1) % lastRound.stitch_count;
        indices.push(
          bottomIdx,
          lastRoundStart + i,
          lastRoundStart + next
        );
      }
    }

    // Create geometry
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute(vertices, 3));
    geometry.setAttribute('normal', new THREE.Float32BufferAttribute(normals, 3));
    geometry.setAttribute('uv', new THREE.Float32BufferAttribute(uvs, 2));
    geometry.setIndex(indices);

    // Compute normals for smooth shading
    geometry.computeVertexNormals();

    return geometry;
  }
}
