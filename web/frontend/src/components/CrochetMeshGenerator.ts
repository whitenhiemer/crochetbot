import * as THREE from 'three';
import { Pattern, Part } from '../types';

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

    // Use pattern-based generation to show what the pattern will actually crochet into
    const part = pattern.parts[0];

    // Determine part type and generate appropriate mesh
    if (part.type === 'sphere') {
      return this.generateSphereMesh(part);
    } else if (part.type === 'cylinder') {
      return this.generateCylinderMesh(part);
    }

    // Analyze pattern to infer shape if type not specified
    const inferredType = this.inferShapeFromPattern(part);
    if (inferredType === 'cylinder') {
      return this.generateCylinderMesh(part);
    }

    // Default to sphere if unknown type
    return this.generateSphereMesh(part);
  }

  /**
   * Generate mesh directly from visualization profile (unsmoothed, accurate)
   */
  private static generateFromVisualizationProfile(profile: number[]): THREE.BufferGeometry {
    const vertices: number[] = [];
    const indices: number[] = [];
    const normals: number[] = [];
    const uvs: number[] = [];

    const numRounds = profile.length;
    const verticesPerRound = 64; // High resolution circles
    const heightScale = 2.5;
    const halfHeight = heightScale / 2;

    // Add bottom center
    vertices.push(0, -halfHeight, 0);
    normals.push(0, -1, 0);
    uvs.push(0.5, 0);

    let vertexIndex = 1;

    // Generate vertices for each profile slice
    for (let roundIdx = 0; roundIdx < numRounds; roundIdx++) {
      const normalizedRadius = profile[roundIdx];
      const radius = normalizedRadius * 0.5;
      const yPos = -halfHeight + (roundIdx / (numRounds - 1)) * heightScale;

      for (let i = 0; i < verticesPerRound; i++) {
        const angle = (i / verticesPerRound) * Math.PI * 2;
        const x = Math.cos(angle) * radius;
        const z = Math.sin(angle) * radius;

        vertices.push(x, yPos, z);

        const normal = new THREE.Vector3(x, 0, z).normalize();
        normals.push(normal.x, normal.y, normal.z);

        uvs.push(i / verticesPerRound, roundIdx / (numRounds - 1));
      }

      // Create faces
      if (roundIdx === 0) {
        // Bottom cap
        for (let i = 0; i < verticesPerRound; i++) {
          const next = (i + 1) % verticesPerRound;
          indices.push(0, 1 + i, 1 + next);
        }
      } else {
        // Side faces
        const prevStart = vertexIndex - verticesPerRound;
        const currStart = vertexIndex;

        for (let i = 0; i < verticesPerRound; i++) {
          const next = (i + 1) % verticesPerRound;

          indices.push(prevStart + i, currStart + next, currStart + i);
          indices.push(prevStart + i, prevStart + next, currStart + next);
        }
      }

      vertexIndex += verticesPerRound;
    }

    // Top cap
    vertices.push(0, halfHeight, 0);
    normals.push(0, 1, 0);
    uvs.push(0.5, 1);
    const topIdx = vertexIndex;
    const lastRoundStart = vertexIndex - verticesPerRound;

    for (let i = 0; i < verticesPerRound; i++) {
      const next = (i + 1) % verticesPerRound;
      indices.push(topIdx, lastRoundStart + next, lastRoundStart + i);
    }

    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute(vertices, 3));
    geometry.setAttribute('normal', new THREE.Float32BufferAttribute(normals, 3));
    geometry.setAttribute('uv', new THREE.Float32BufferAttribute(uvs, 2));
    geometry.setIndex(indices);

    geometry.computeVertexNormals();

    return geometry;
  }

  /**
   * Infer shape type from pattern structure
   */
  private static inferShapeFromPattern(part: Part): 'sphere' | 'cylinder' {
    const rounds = part.rounds.filter(r => r.stitch_count > 0);
    if (rounds.length === 0) return 'sphere';

    // Count constant rounds
    let constants = 0;

    for (let i = 1; i < rounds.length; i++) {
      const diff = rounds[i].stitch_count - rounds[i - 1].stitch_count;
      if (diff === 0) constants++;
    }

    // If mostly constant rounds with few increases/decreases at ends, it's a cylinder
    const constantRatio = constants / rounds.length;
    if (constantRatio > 0.5) {
      return 'cylinder';
    }

    return 'sphere';
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

      // Calculate radius at this round based on actual stitch count
      const maxStitches = Math.max(...rounds.map(r => r.stitch_count));
      const normalizedHeight = (roundIdx + 1) / totalRounds;

      // Amplify stitch count differences for visible shape
      const stitchRatio = stitchCount / maxStitches;
      const amplifiedRatio = 0.3 + (stitchRatio * 0.7); // Map 0-1 to 0.3-1.0
      let radius = amplifiedRatio * 0.6;

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

  /**
   * Generate a cylindrical crocheted mesh from rounds
   */
  private static generateCylinderMesh(part: Part): THREE.BufferGeometry {
    const rounds = part.rounds.filter(r => r.stitch_count > 0);

    if (rounds.length === 0) {
      return new THREE.CylinderGeometry(0.5, 0.5, 1, 12);
    }

    const vertices: number[] = [];
    const indices: number[] = [];
    const normals: number[] = [];
    const uvs: number[] = [];

    const maxStitches = Math.max(...rounds.map(r => r.stitch_count));

    // Scale total height based on round count (more rounds = taller)
    const heightScale = Math.min(3.0, 1.0 + (rounds.length / 30));
    const totalHeight = heightScale;
    const halfHeight = totalHeight / 2;

    // Add bottom center point
    vertices.push(0, -halfHeight, 0);
    normals.push(0, -1, 0);
    uvs.push(0.5, 0);

    let vertexIndex = 1;
    const roundVertexCounts: number[] = []; // Track vertices per round for face generation

    // Generate vertices for each round
    for (let roundIdx = 0; roundIdx < rounds.length; roundIdx++) {
      const round = rounds[roundIdx];
      const stitchCount = round.stitch_count;

      if (stitchCount === 0) continue;

      // Calculate Y position - evenly space all rounds from bottom to top
      const yPos = -halfHeight + (roundIdx / (rounds.length - 1)) * totalHeight;

      // Calculate radius at this round based on stitch count
      // Amplify the stitch count differences to make shape variation more visible
      const stitchRatio = stitchCount / maxStitches;
      const amplifiedRatio = 0.3 + (stitchRatio * 0.7); // Map 0-1 to 0.3-1.0 (ensure minimum width)
      let radius = amplifiedRatio * 0.6; // Slightly larger base radius

      // Add slight variation for hand-crocheted look
      const bumpiness = 0.04;

      // Create more vertices than stitches for smoother appearance
      // Use at least 48 vertices per round for very smooth circles
      const minVertices = 48;
      const verticesThisRound = Math.max(minVertices, stitchCount * 3);

      // Create vertices around this round
      for (let i = 0; i < verticesThisRound; i++) {
        const angle = (i / verticesThisRound) * Math.PI * 2;

        // Add slight bumps at stitch positions
        const stitchAngle = (Math.floor((i / verticesThisRound) * stitchCount) / stitchCount) * Math.PI * 2;
        const stitchBump = Math.sin((angle - stitchAngle) * stitchCount) * bumpiness * 0.5;

        const radiusVariation = radius + stitchBump;

        const x = Math.cos(angle) * radiusVariation;
        const z = Math.sin(angle) * radiusVariation;
        const y = yPos + Math.sin(i * 3) * bumpiness * 0.3;

        vertices.push(x, y, z);

        // Calculate normal (pointing outward)
        const normal = new THREE.Vector3(x, 0, z).normalize();
        normals.push(normal.x, normal.y, normal.z);

        // UV coordinates
        uvs.push(i / verticesThisRound, (roundIdx + 1) / rounds.length);
      }

      // Store vertex count for this round
      roundVertexCounts.push(verticesThisRound);

      // Create faces connecting to previous round
      if (roundIdx === 0) {
        // Connect first round to bottom center
        const currentRoundStart = 1;
        for (let i = 0; i < verticesThisRound; i++) {
          const next = (i + 1) % verticesThisRound;
          indices.push(
            0,
            currentRoundStart + i,
            currentRoundStart + next
          );
        }
      } else {
        // Connect this round to previous round
        const prevVertexCount = roundVertexCounts[roundIdx - 1];
        const currVertexCount = verticesThisRound;
        const prevRoundStart = vertexIndex - prevVertexCount;
        const currentRoundStart = vertexIndex;

        const ratio = currVertexCount / prevVertexCount;

        if (ratio >= 1) {
          // Increasing or same
          for (let i = 0; i < prevVertexCount; i++) {
            const prevIdx = prevRoundStart + i;
            const prevNext = prevRoundStart + ((i + 1) % prevVertexCount);

            const currStart = Math.floor(i * ratio);
            const currEnd = Math.floor((i + 1) * ratio);

            for (let j = currStart; j < currEnd; j++) {
              const currIdx = currentRoundStart + (j % currVertexCount);
              const currNext = currentRoundStart + ((j + 1) % currVertexCount);

              indices.push(prevIdx, currNext, currIdx);
              indices.push(prevIdx, prevNext, currNext);
            }
          }
        } else {
          // Decreasing
          for (let i = 0; i < currVertexCount; i++) {
            const currIdx = currentRoundStart + i;
            const currNext = currentRoundStart + ((i + 1) % currVertexCount);

            const prevStart = Math.floor(i / ratio);
            const prevEnd = Math.floor((i + 1) / ratio);

            for (let j = prevStart; j <= prevEnd && j < prevVertexCount; j++) {
              const prevIdx = prevRoundStart + (j % prevVertexCount);
              const prevNext = prevRoundStart + ((j + 1) % prevVertexCount);

              indices.push(prevIdx, currNext, currIdx);
              if (j < prevEnd) {
                indices.push(prevIdx, prevNext, currNext);
              }
            }
          }
        }
      }

      vertexIndex += verticesThisRound;
    }

    // Add top center point
    if (roundVertexCounts.length > 0) {
      const lastVertexCount = roundVertexCounts[roundVertexCounts.length - 1];
      const lastRoundStart = vertexIndex - lastVertexCount;

      vertices.push(0, halfHeight, 0);
      normals.push(0, 1, 0);
      uvs.push(0.5, 1);

      const topIdx = vertexIndex;

      // Connect last round to top center
      for (let i = 0; i < lastVertexCount; i++) {
        const next = (i + 1) % lastVertexCount;
        indices.push(
          topIdx,
          lastRoundStart + next,
          lastRoundStart + i
        );
      }
    }

    // Create geometry
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute(vertices, 3));
    geometry.setAttribute('normal', new THREE.Float32BufferAttribute(normals, 3));
    geometry.setAttribute('uv', new THREE.Float32BufferAttribute(uvs, 2));
    geometry.setIndex(indices);

    geometry.computeVertexNormals();

    return geometry;
  }
}
