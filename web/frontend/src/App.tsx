import React, { useState } from 'react';
import { FileUpload } from './components/FileUpload';
import { ModelPreview } from './components/ModelPreview';
import { PatternDisplay } from './components/PatternDisplay';
import { api } from './api';
import { AppStep, MeshFile, Pattern } from './types';
import './App.css';

function App() {
  const [step, setStep] = useState<AppStep>('upload');
  const [file, setFile] = useState<File | null>(null);
  const [meshData, setMeshData] = useState<MeshFile | null>(null);
  const [pattern, setPattern] = useState<Pattern | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileSelected = async (selectedFile: File) => {
    setFile(selectedFile);
    setError(null);
    setLoading(true);

    try {
      const response = await api.uploadFile(selectedFile);

      if (response.success && response.file) {
        setMeshData(response.file);
        setStep('preview');
      } else {
        setError(response.error || 'Upload failed');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to upload file');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerate = async () => {
    if (!meshData) return;

    setError(null);
    setLoading(true);

    try {
      const response = await api.generatePattern(meshData.filename, meshData.id);

      if (response.success && response.pattern) {
        setPattern(response.pattern);
        setStep('pattern');
      } else {
        setError(response.error || 'Pattern generation failed');
      }
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to generate pattern');
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setStep('upload');
    setFile(null);
    setMeshData(null);
    setPattern(null);
    setError(null);
  };

  return (
    <div className="App">
      {error && (
        <div className="error-banner">
          <span className="error-icon">⚠️</span>
          <span>{error}</span>
          <button className="error-close" onClick={() => setError(null)}>
            ×
          </button>
        </div>
      )}

      {step === 'upload' && (
        <FileUpload onFileSelected={handleFileSelected} loading={loading} />
      )}

      {step === 'preview' && file && (
        <ModelPreview
          file={file}
          meshData={meshData}
          onGenerate={handleGenerate}
          generating={loading}
        />
      )}

      {step === 'pattern' && pattern && file && (
        <PatternDisplay pattern={pattern} file={file} onReset={handleReset} />
      )}

      <footer className="app-footer">
        <p>
          CrochetBot - Convert 3D models to amigurumi patterns
        </p>
      </footer>
    </div>
  );
}

export default App;
