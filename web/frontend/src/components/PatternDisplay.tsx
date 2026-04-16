import React, { useState, useEffect } from 'react';
import { Pattern } from '../types';
import { YarnPreview } from './YarnPreview';
import './PatternDisplay.css';

interface PatternDisplayProps {
  pattern: Pattern;
  file: File;
  onReset: () => void;
}

export const PatternDisplay: React.FC<PatternDisplayProps> = ({ pattern, file, onReset }) => {
  const downloadPattern = () => {
    const patternText = generatePatternText(pattern);
    const blob = new Blob([patternText], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${pattern.name.replace(/\s+/g, '_')}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const downloadJSON = () => {
    const blob = new Blob([JSON.stringify(pattern, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${pattern.name.replace(/\s+/g, '_')}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="pattern-display-container">
      <div className="pattern-header">
        <div>
          <h1>{pattern.name}</h1>
          <p className="pattern-meta">
            Difficulty: <span className={`difficulty ${pattern.difficulty}`}>{pattern.difficulty}</span>
          </p>
        </div>
        <div className="header-actions">
          <button onClick={downloadPattern} className="download-button">
            Download Pattern
          </button>
          <button onClick={downloadJSON} className="download-button secondary">
            Download JSON
          </button>
          <button onClick={onReset} className="reset-button">
            New Pattern
          </button>
        </div>
      </div>

      <div className="pattern-content">
        {/* Yarn Preview Section */}
        <YarnPreview pattern={pattern} />

        {/* Materials Section */}
        <section className="pattern-section materials-section">
          <h2>Materials</h2>
          <div className="materials-grid">
            <div className="material-item">
              <strong>Yarn:</strong> {pattern.materials.yarn_weight} weight (~
              {pattern.materials.yarn_yardage} yards)
            </div>
            <div className="material-item">
              <strong>Hook:</strong> {pattern.materials.hook_size}
            </div>
            {pattern.materials.colors.map((color, idx) => (
              <div key={idx} className="material-item">
                <strong>Color:</strong> {color.name} ({color.amount} yards)
              </div>
            ))}
          </div>
          {pattern.materials.other_supplies.length > 0 && (
            <div className="supplies">
              <strong>Other Supplies:</strong>
              <ul>
                {pattern.materials.other_supplies.map((supply, idx) => (
                  <li key={idx}>{supply}</li>
                ))}
              </ul>
            </div>
          )}
        </section>

        {/* Parts and Rounds */}
        {pattern.parts.map((part, partIdx) => (
          <section key={partIdx} className="pattern-section part-section">
            <h2>
              {part.name} <span className="part-type">({part.type})</span>
            </h2>
            <p className="starting-info">
              <strong>Starting:</strong> {part.starting_type}
            </p>

            <div className="rounds-list">
              {part.rounds.map((round, roundIdx) => (
                <div key={roundIdx} className="round-item">
                  <div className="round-header">
                    <span className="round-number">Round {round.number}</span>
                    <span className="stitch-count">{round.stitch_count} sts</span>
                  </div>
                  <div className="round-instructions">{round.instructions}</div>
                  {round.notes && (
                    <div className="round-notes">
                      <span className="note-icon">💡</span>
                      {round.notes}
                    </div>
                  )}
                </div>
              ))}
            </div>

            {part.notes.length > 0 && (
              <div className="part-notes">
                <strong>Notes:</strong>
                <ul>
                  {part.notes.map((note, idx) => (
                    <li key={idx}>{note}</li>
                  ))}
                </ul>
              </div>
            )}
          </section>
        ))}

        {/* Assembly Instructions */}
        {pattern.assembly_instructions.length > 0 && (
          <section className="pattern-section assembly-section">
            <h2>Assembly</h2>
            <ol>
              {pattern.assembly_instructions.map((instruction, idx) => (
                <li key={idx}>{instruction}</li>
              ))}
            </ol>
          </section>
        )}
      </div>
    </div>
  );
};

// Helper function to generate plain text pattern
function generatePatternText(pattern: Pattern): string {
  let text = `${pattern.name}\n`;
  text += `${'='.repeat(pattern.name.length)}\n\n`;
  text += `${pattern.description}\n`;
  text += `Difficulty: ${pattern.difficulty}\n\n`;

  text += `MATERIALS\n---------\n`;
  text += `- Yarn: ${pattern.materials.yarn_weight} weight, ~${pattern.materials.yarn_yardage} yards\n`;
  text += `- Hook: ${pattern.materials.hook_size}\n`;
  pattern.materials.colors.forEach((color) => {
    text += `- ${color.name}: ${color.amount} yards\n`;
  });
  pattern.materials.other_supplies.forEach((supply) => {
    text += `- ${supply}\n`;
  });
  text += `\n`;

  pattern.parts.forEach((part) => {
    text += `${part.name.toUpperCase()} (${part.type})\n`;
    text += `${'-'.repeat(part.name.length + part.type.length + 3)}\n`;
    text += `Starting: ${part.starting_type}\n\n`;

    part.rounds.forEach((round) => {
      text += `Round ${round.number}: ${round.instructions} (${round.stitch_count} sts)\n`;
      if (round.notes) {
        text += `  Note: ${round.notes}\n`;
      }
    });

    if (part.notes.length > 0) {
      text += `\nNotes:\n`;
      part.notes.forEach((note) => {
        text += `- ${note}\n`;
      });
    }
    text += `\n`;
  });

  if (pattern.assembly_instructions.length > 0) {
    text += `ASSEMBLY\n--------\n`;
    pattern.assembly_instructions.forEach((instruction, idx) => {
      text += `${idx + 1}. ${instruction}\n`;
    });
  }

  return text;
}
