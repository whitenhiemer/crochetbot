import axios from 'axios';
import { UploadResponse, GenerateResponse, Pattern } from './types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const api = {
  // Health check
  async health(): Promise<{ status: string; service: string }> {
    const response = await axios.get(`${API_BASE_URL}/health`);
    return response.data;
  },

  // Upload file
  async uploadFile(file: File): Promise<UploadResponse> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await axios.post(`${API_BASE_URL}/api/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    return response.data;
  },

  // Generate pattern
  async generatePattern(filename: string, fileId?: string): Promise<GenerateResponse> {
    const response = await axios.post(`${API_BASE_URL}/api/generate`, {
      filename,
      file_id: fileId,
    });

    return response.data;
  },

  // Get pattern by ID
  async getPattern(patternId: string): Promise<Pattern> {
    const response = await axios.get(`${API_BASE_URL}/api/pattern/${patternId}`);
    return response.data;
  },
};
