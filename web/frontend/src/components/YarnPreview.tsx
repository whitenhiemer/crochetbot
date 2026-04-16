import React, { useEffect, useState } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Center, Environment } from '@react-three/drei';
import * as THREE from 'three';
import { OBJLoader } from 'three/examples/jsm/loaders/OBJLoader.js';
import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
import './YarnPreview.css';

interface YarnPreviewProps {
  fileUrl: string;
  fileType: string;
  yarnColor?: string;
}

const YarnMesh: React.FC<{ fileUrl: string; fileType: string; yarnColor: string }> = ({
  fileUrl,
  fileType,
  yarnColor
}) => {
  const [geometry, setGeometry] = useState<THREE.BufferGeometry | null>(null);

  useEffect(() => {
    if (fileType === 'obj') {
      const loader = new OBJLoader();
      loader.load(fileUrl, (object) => {
        const mesh = object.children.find((child) => child instanceof THREE.Mesh) as THREE.Mesh;
        if (mesh && mesh.geometry) {
          setGeometry(mesh.geometry);
        }
      });
    } else if (fileType === 'stl') {
      const loader = new STLLoader();
      loader.load(fileUrl, (geometry) => {
        setGeometry(geometry);
      });
    }
  }, [fileUrl, fileType]);

  if (!geometry) {
    return null;
  }

  // Create yarn-like texture
  const textureCanvas = document.createElement('canvas');
  textureCanvas.width = 512;
  textureCanvas.height = 512;
  const ctx = textureCanvas.getContext('2d');

  if (ctx) {
    // Base color
    ctx.fillStyle = yarnColor;
    ctx.fillRect(0, 0, 512, 512);

    // Add yarn texture pattern (diagonal lines to simulate stitches)
    ctx.strokeStyle = adjustBrightness(yarnColor, -20);
    ctx.lineWidth = 2;

    for (let i = 0; i < 512; i += 8) {
      ctx.beginPath();
      ctx.moveTo(i, 0);
      ctx.lineTo(i + 20, 512);
      ctx.stroke();

      ctx.beginPath();
      ctx.moveTo(0, i);
      ctx.lineTo(512, i + 20);
      ctx.stroke();
    }

    // Add highlights (yarn sheen)
    ctx.strokeStyle = adjustBrightness(yarnColor, 40);
    ctx.lineWidth = 1;

    for (let i = 0; i < 512; i += 16) {
      ctx.beginPath();
      ctx.moveTo(i + 2, 0);
      ctx.lineTo(i + 22, 512);
      ctx.stroke();
    }
  }

  const texture = new THREE.CanvasTexture(textureCanvas);
  texture.wrapS = THREE.RepeatWrapping;
  texture.wrapT = THREE.RepeatWrapping;
  texture.repeat.set(4, 4);

  return (
    <Center>
      <mesh geometry={geometry}>
        <meshStandardMaterial
          map={texture}
          color={yarnColor}
          roughness={0.8}
          metalness={0.1}
          bumpMap={texture}
          bumpScale={0.02}
        />
      </mesh>
    </Center>
  );
};

// Helper function to adjust color brightness
function adjustBrightness(color: string, amount: number): string {
  const hex = color.replace('#', '');
  const r = Math.max(0, Math.min(255, parseInt(hex.substr(0, 2), 16) + amount));
  const g = Math.max(0, Math.min(255, parseInt(hex.substr(2, 2), 16) + amount));
  const b = Math.max(0, Math.min(255, parseInt(hex.substr(4, 2), 16) + amount));

  return '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('');
}

export const YarnPreview: React.FC<YarnPreviewProps> = ({
  fileUrl,
  fileType,
  yarnColor = '#e91e63'
}) => {
  const [selectedColor, setSelectedColor] = useState(yarnColor);

  const yarnColors = [
    { name: 'Pink', value: '#e91e63' },
    { name: 'Blue', value: '#2196f3' },
    { name: 'Green', value: '#4caf50' },
    { name: 'Yellow', value: '#ffeb3b' },
    { name: 'Red', value: '#f44336' },
    { name: 'Purple', value: '#9c27b0' },
    { name: 'Orange', value: '#ff9800' },
    { name: 'Brown', value: '#795548' },
    { name: 'Gray', value: '#9e9e9e' },
    { name: 'White', value: '#f5f5f5' },
    { name: 'Black', value: '#212121' },
    { name: 'Teal', value: '#009688' },
  ];

  return (
    <div className="yarn-preview-container">
      <div className="yarn-preview-header">
        <h3>Yarn Preview</h3>
        <p>See how your pattern might look in yarn</p>
      </div>

      <div className="yarn-canvas-container">
        <Canvas camera={{ position: [2, 2, 2], fov: 50 }}>
          <ambientLight intensity={0.6} />
          <directionalLight position={[10, 10, 5]} intensity={0.8} />
          <spotLight position={[-10, 10, -5]} intensity={0.3} />
          <YarnMesh fileUrl={fileUrl} fileType={fileType} yarnColor={selectedColor} />
          <OrbitControls enablePan={true} enableZoom={true} enableRotate={true} />
          <Environment preset="apartment" />
        </Canvas>
      </div>

      <div className="color-picker">
        <label>Yarn Color:</label>
        <div className="color-swatches">
          {yarnColors.map((color) => (
            <button
              key={color.value}
              className={`color-swatch ${selectedColor === color.value ? 'selected' : ''}`}
              style={{ backgroundColor: color.value }}
              onClick={() => setSelectedColor(color.value)}
              title={color.name}
              aria-label={`Select ${color.name} yarn color`}
            />
          ))}
        </div>
      </div>

      <div className="preview-note">
        <p>
          <strong>Note:</strong> This is an approximation. Actual crocheted result may vary based on
          yarn weight, hook size, tension, and skill level.
        </p>
      </div>
    </div>
  );
};
