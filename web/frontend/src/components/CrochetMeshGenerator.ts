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
    const vertexCounts: number[] = [];

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

      // Stuffed sphere: fabric stretches and rounds out
      const stitchRatio = stitchCount / maxStitches;
      const stuffingFactor = 1.2; // 20% inflation from stuffing

      // Light smoothing for stuffed appearance while preserving detail
      const smoothingWindow = 1;
      let avgRatio = stitchRatio;
      if (roundIdx >= smoothingWindow && roundIdx < rounds.length - smoothingWindow) {
        let sum = 0;
        let count = 0;
        for (let j = -smoothingWindow; j <= smoothingWindow; j++) {
          if (roundIdx + j >= 0 && roundIdx + j < rounds.length) {
            sum += (rounds[roundIdx + j].stitch_count / maxStitches);
            count++;
          }
        }
        avgRatio = sum / count;
      }

      let radius = avgRatio * stuffingFactor * 0.85;

      // Minimal texture - stuffed pieces are smooth
      const bumpiness = 0.01;

      // High vertex count for smooth appearance
      const verticesThisRound = Math.max(64, stitchCount * 4);

      // Create vertices around this round
      for (let i = 0; i < verticesThisRound; i++) {
        const angle = (i / verticesThisRound) * Math.PI * 2;

        // Very subtle texture
        const radiusVariation = radius + Math.sin(i * 6.3) * bumpiness;

        const x = Math.cos(angle) * radiusVariation;
        const z = Math.sin(angle) * radiusVariation;
        const y = yPos;

        vertices.push(x, y, z);

        // Calculate normal (pointing outward)
        const normal = new THREE.Vector3(x, y, z).normalize();
        normals.push(normal.x, normal.y, normal.z);

        // UV coordinates
        uvs.push(i / stitchCount, 1 - normalizedHeight);
      }

      // Track vertex counts per round
      if (roundIdx === 0) {
        vertexCounts.push(verticesThisRound);
      } else {
        vertexCounts.push(verticesThisRound);
      }

      // Create faces connecting to previous round
      if (roundIdx === 0) {
        // Connect first round to center point
        const currentRoundStart = 1;
        for (let i = 0; i < verticesThisRound; i++) {
          const next = (i + 1) % verticesThisRound;
          indices.push(
            0, // Center point
            currentRoundStart + next,
            currentRoundStart + i
          );
        }
      } else {
        // Connect this round to previous round
        const prevVertexCount = vertexCounts[roundIdx - 1];
        const currVertexCount = verticesThisRound;
        const prevRoundStart = vertexIndex - prevVertexCount;
        const currentRoundStart = vertexIndex;

        // Create faces between rounds with different vertex counts
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

    // Add bottom center point if last round has vertices
    if (vertexCounts.length > 0) {
      const lastVertexCount = vertexCounts[vertexCounts.length - 1];
      const lastRoundStart = vertexIndex - lastVertexCount;

      vertices.push(0, -0.5, 0);
      normals.push(0, -1, 0);
      uvs.push(0.5, 0);

      const bottomIdx = vertexIndex;

      // Connect last round to bottom center
      for (let i = 0; i < lastVertexCount; i++) {
        const next = (i + 1) % lastVertexCount;
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
      // Stuffed amigurumi: fabric stretches and rounds out, making it fuller than unstuffed
      const stitchRatio = stitchCount / maxStitches;

      // Apply stuffing effect: smaller stitch counts get more inflated when stuffed
      // This creates rounder, more organic shapes
      const stuffingFactor = 1.15; // 15% inflation from stuffing
      const baseRadius = stitchRatio * stuffingFactor;

      // Light smoothing - preserve detail while showing stuffed roundness
      const smoothingWindow = 1;
      let avgRatio = baseRadius;
      if (roundIdx >= smoothingWindow && roundIdx < rounds.length - smoothingWindow) {
        let sum = 0;
        let count = 0;
        for (let j = -smoothingWindow; j <= smoothingWindow; j++) {
          if (roundIdx + j >= 0 && roundIdx + j < rounds.length) {
            sum += (rounds[roundIdx + j].stitch_count / maxStitches);
            count++;
          }
        }
        avgRatio = (sum / count) * stuffingFactor;
      }

      let radius = avgRatio * 0.9;

      // Subtle texture only
      const bumpiness = 0.02;

      // High vertex count for smooth stuffed appearance
      const minVertices = 64;
      const verticesThisRound = Math.max(minVertices, stitchCount * 4);

      // Create vertices around this round
      for (let i = 0; i < verticesThisRound; i++) {
        const angle = (i / verticesThisRound) * Math.PI * 2;

        // Very subtle texture - stuffed pieces are smooth
        const textureNoise = Math.sin(i * 5.7) * bumpiness;

        const radiusVariation = radius + textureNoise;

        const x = Math.cos(angle) * radiusVariation;
        const z = Math.sin(angle) * radiusVariation;
        const y = yPos;

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
