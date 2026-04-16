import React, { useEffect, useState } from 'react';
import { Canvas, useThree } from '@react-three/fiber';
import { OrbitControls, Center, Environment } from '@react-three/drei';
import * as THREE from 'three';
import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader.js';
import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
import { MeshFile } from '../types';
import './ModelPreview.css';

interface ModelPreviewProps {
  file: File;
  meshData: MeshFile | null;
  onGenerate: () => void;
  generating: boolean;
}

const CameraController: React.FC<{ geometry: THREE.BufferGeometry | null }> = ({ geometry }) => {
  const { camera } = useThree();

  useEffect(() => {
    if (!geometry) return;

    // Compute bounding box
    geometry.computeBoundingBox();
    const boundingBox = geometry.boundingBox;

    if (!boundingBox) return;

    // Get center and size
    const center = new THREE.Vector3();
    boundingBox.getCenter(center);
    const size = new THREE.Vector3();
    boundingBox.getSize(size);

    // Get the max dimension
    const maxDim = Math.max(size.x, size.y, size.z);

    // Calculate camera distance (with some padding)
    const fov = (camera as THREE.PerspectiveCamera).fov * (Math.PI / 180);
    const cameraDistance = Math.abs(maxDim / Math.sin(fov / 2)) * 1.5;

    // Determine optimal camera position based on model orientation
    let cameraPos = new THREE.Vector3();

    // Analyze model aspect ratio to determine orientation
    const aspectXY = size.x / size.y;
    const aspectYZ = size.y / size.z;

    // Determine dominant axis and orientation
    if (size.y > size.x && size.y > size.z && aspectYZ > 2) {
      // Tall/vertical model (like a person standing)
      // View from slightly elevated front-right angle
      cameraPos.set(
        center.x + cameraDistance * 0.6,
        center.y + cameraDistance * 0.3,
        center.z + cameraDistance * 0.7
      );
    } else if (size.z > size.x && size.z > size.y && aspectYZ < 0.5) {
      // Deep model (extending along Z)
      // View from elevated side angle
      cameraPos.set(
        center.x + cameraDistance * 0.7,
        center.y + cameraDistance * 0.5,
        center.z + cameraDistance * 0.5
      );
    } else if (size.x > size.y && size.x > size.z && aspectXY > 2) {
      // Wide/horizontal model
      // View from front-elevated angle
      cameraPos.set(
        center.x + cameraDistance * 0.3,
        center.y + cameraDistance * 0.5,
        center.z + cameraDistance * 0.8
      );
    } else if (size.y < maxDim * 0.3) {
      // Flat/low model (like a coin or base)
      // View from above at 45 degree angle
      cameraPos.set(
        center.x + cameraDistance * 0.5,
        center.y + cameraDistance * 0.9,
        center.z + cameraDistance * 0.5
      );
    } else {
      // Roughly cubic/spherical model
      // Classic 3/4 view angle
      cameraPos.set(
        center.x + cameraDistance * 0.6,
        center.y + cameraDistance * 0.6,
        center.z + cameraDistance * 0.6
      );
    }

    camera.position.copy(cameraPos);
    camera.lookAt(center);
    camera.updateProjectionMatrix();
  }, [geometry, camera]);

  return null;
};

const Model: React.FC<{ fileUrl: string; fileType: string }> = ({ fileUrl, fileType }) => {
  const [geometry, setGeometry] = useState<THREE.BufferGeometry | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setGeometry(null);
    setError(null);

    if (fileType === 'obj') {
      const loader = new OBJLoader();
      loader.load(
        fileUrl,
        (object) => {
          // Extract geometry from loaded object
          const mesh = object.children.find((child) => child instanceof THREE.Mesh) as THREE.Mesh;
          if (mesh && mesh.geometry) {
            setGeometry(mesh.geometry);
          } else {
            setError('No mesh found in OBJ file');
          }
        },
        undefined,
        (error) => {
          console.error('Error loading OBJ:', error);
          setError('Failed to load OBJ file');
        }
      );
    } else if (fileType === 'stl') {
      const loader = new STLLoader();
      loader.load(
        fileUrl,
        (geometry) => {
          setGeometry(geometry);
        },
        undefined,
        (error) => {
          console.error('Error loading STL:', error);
          setError('Failed to load STL file');
        }
      );
    }
  }, [fileUrl, fileType]);

  if (error) {
    return (
      <Center>
        <mesh>
          <boxGeometry args={[1, 1, 1]} />
          <meshStandardMaterial color="#ff0000" />
        </mesh>
      </Center>
    );
  }

  if (!geometry) {
    return null;
  }

  return (
    <>
      <CameraController geometry={geometry} />
      <Center>
        <mesh geometry={geometry}>
          <meshStandardMaterial color="#e91e63" />
        </mesh>
      </Center>
    </>
  );
};

export const ModelPreview: React.FC<ModelPreviewProps> = ({
  file,
  meshData,
  onGenerate,
  generating,
}) => {
  const [fileUrl, setFileUrl] = useState<string | null>(null);
  const [fileType, setFileType] = useState<string>('');

  useEffect(() => {
    // Create object URL for the file
    const url = URL.createObjectURL(file);
    setFileUrl(url);

    // Determine file type from extension
    const extension = file.name.split('.').pop()?.toLowerCase() || '';
    setFileType(extension);

    // Cleanup
    return () => {
      URL.revokeObjectURL(url);
    };
  }, [file]);

  return (
    <div className="model-preview-container">
      <h2>3D Model Preview</h2>

      <div className="preview-card">
        <div className="canvas-container">
          {fileUrl && fileType ? (
            <Canvas camera={{ position: [2, 2, 2], fov: 50 }}>
              <ambientLight intensity={0.5} />
              <directionalLight position={[10, 10, 5]} intensity={1} />
              <Model fileUrl={fileUrl} fileType={fileType} />
              <OrbitControls enablePan={true} enableZoom={true} enableRotate={true} />
              <Environment preset="studio" />
            </Canvas>
          ) : (
            <div className="loading-preview">Loading model...</div>
          )}
        </div>

        <div className="model-info">
          <h3>Model Information</h3>
          <div className="info-grid">
            <div className="info-item">
              <span className="info-label">File:</span>
              <span className="info-value">{file.name}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Size:</span>
              <span className="info-value">{(file.size / 1024).toFixed(2)} KB</span>
            </div>
            <div className="info-item">
              <span className="info-label">Format:</span>
              <span className="info-value">{meshData?.format.toUpperCase() || 'OBJ'}</span>
            </div>
            {meshData && meshData.vertices > 0 && (
              <>
                <div className="info-item">
                  <span className="info-label">Vertices:</span>
                  <span className="info-value">{meshData.vertices}</span>
                </div>
                <div className="info-item">
                  <span className="info-label">Faces:</span>
                  <span className="info-value">{meshData.faces}</span>
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      <div className="action-buttons">
        <button
          className="generate-button"
          onClick={onGenerate}
          disabled={generating}
        >
          {generating ? 'Generating Pattern...' : 'Generate Crochet Pattern'}
        </button>
      </div>

      <div className="instructions">
        <p>
          <strong>Preview Controls:</strong> Click and drag to rotate, scroll to zoom, right-click to pan
        </p>
      </div>
    </div>
  );
};
