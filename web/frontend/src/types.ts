// API Types
export interface UploadResponse {
  success: boolean;
  message: string;
  file?: MeshFile;
  error?: string;
}

export interface MeshFile {
  id: string;
  filename: string;
  uploaded_at: string;
  vertices: number;
  faces: number;
  format: string;
}

export interface GenerateResponse {
  success: boolean;
  message: string;
  pattern?: Pattern;
  error?: string;
}

export interface Pattern {
  id: string;
  name: string;
  created_at: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  parts: Part[];
  materials: Materials;
  assembly_instructions: string[];
  finished_size: FinishedSize;
  accuracy_metrics: AccuracyMetrics;
  visualization_profile?: number[]; // Raw mesh radius profile (0-1) for accurate rendering
}

export interface FinishedSize {
  height_inches: number;
  height_cm: number;
  width_inches: number;
  width_cm: number;
}

export interface AccuracyMetrics {
  shape_match_percent: number;
  average_error: number;
  max_error: number;
  notes: string;
}

export interface Part {
  name: string;
  type: string;
  rounds: Round[];
  color: string;
  starting_type: string;
  notes: string[];
}

export interface Round {
  number: number;
  instructions: string;
  stitch_count: number;
  stitch_type: string;
  repeats: number;
  notes: string;
}

export interface Materials {
  yarn_weight: string;
  yarn_yardage: number;
  hook_size: string;
  colors: Color[];
  other_supplies: string[];
}

export interface Color {
  name: string;
  amount: number;
}

// UI State
export type AppStep = 'upload' | 'preview' | 'pattern';

export interface AppState {
  step: AppStep;
  file: File | null;
  meshData: MeshFile | null;
  pattern: Pattern | null;
  loading: boolean;
  error: string | null;
}
