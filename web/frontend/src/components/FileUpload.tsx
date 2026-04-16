import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import './FileUpload.css';

interface FileUploadProps {
  onFileSelected: (file: File) => void;
  loading: boolean;
}

export const FileUpload: React.FC<FileUploadProps> = ({ onFileSelected, loading }) => {
  const onDrop = useCallback(
    (acceptedFiles: File[]) => {
      if (acceptedFiles.length > 0) {
        onFileSelected(acceptedFiles[0]);
      }
    },
    [onFileSelected]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'model/obj': ['.obj'],
      // 'model/stl': ['.stl'], // TODO: STL not yet supported
    },
    multiple: false,
    disabled: loading,
  });

  return (
    <div className="file-upload-container">
      <h1>CrochetBot</h1>
      <p className="subtitle">Convert 3D models to amigurumi crochet patterns</p>

      <div
        {...getRootProps()}
        className={`dropzone ${isDragActive ? 'active' : ''} ${loading ? 'disabled' : ''}`}
      >
        <input {...getInputProps()} />
        <div className="dropzone-content">
          <svg
            className="upload-icon"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
            />
          </svg>
          {loading ? (
            <p>Uploading...</p>
          ) : isDragActive ? (
            <p>Drop the file here</p>
          ) : (
            <>
              <p>Drag and drop a 3D model file here</p>
              <p className="or-text">or</p>
              <button className="browse-button">Browse Files</button>
              <p className="file-types">Supported: .obj (max 50MB)</p>
            </>
          )}
        </div>
      </div>

      <div className="features">
        <div className="feature">
          <div className="feature-icon">📐</div>
          <h3>Upload Model</h3>
          <p>Upload your 3D model file</p>
        </div>
        <div className="feature">
          <div className="feature-icon">🧶</div>
          <h3>Generate Pattern</h3>
          <p>AI converts to crochet instructions</p>
        </div>
        <div className="feature">
          <div className="feature-icon">📄</div>
          <h3>Download</h3>
          <p>Get your complete pattern</p>
        </div>
      </div>
    </div>
  );
};
