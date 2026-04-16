import React, { useEffect, useState } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Center, Environment } from '@react-three/drei';
import * as THREE from 'three';
import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader.js';
import { MeshFile } from '../types';
import './ModelPreview.css';

interface ModelPreviewProps {
  file: File;
  meshData: MeshFile | null;
  onGenerate: () => void;
  generating: boolean;
}

const Model: React.FC<{ objUrl: string }> = ({ objUrl }) => {
  const [geometry, setGeometry] = useState<THREE.BufferGeometry | null>(null);

  useEffect(() => {
    const loader = new OBJLoader();
    loader.load(
      objUrl,
      (object) => {
        // Extract geometry from loaded object
        const mesh = object.children.find((child) => child instanceof THREE.Mesh) as THREE.Mesh;
        if (mesh && mesh.geometry) {
          setGeometry(mesh.geometry);
        }
      },
      undefined,
      (error) => {
        console.error('Error loading OBJ:', error);
      }
    );
  }, [objUrl]);

  if (!geometry) {
    return null;
  }

  return (
    <Center>
      <mesh geometry={geometry}>
        <meshStandardMaterial color="#e91e63" />
      </mesh>
    </Center>
  );
};

export const ModelPreview: React.FC<ModelPreviewProps> = ({
  file,
  meshData,
  onGenerate,
  generating,
}) => {
  const [objUrl, setObjUrl] = useState<string | null>(null);

  useEffect(() => {
    // Create object URL for the file
    const url = URL.createObjectURL(file);
    setObjUrl(url);

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
          {objUrl ? (
            <Canvas camera={{ position: [2, 2, 2], fov: 50 }}>
              <ambientLight intensity={0.5} />
              <directionalLight position={[10, 10, 5]} intensity={1} />
              <Model objUrl={objUrl} />
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
