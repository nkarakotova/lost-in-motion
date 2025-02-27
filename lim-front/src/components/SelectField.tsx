import React from 'react';
import { MenuItem, Select, InputLabel, FormControl, SelectChangeEvent } from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface SelectFieldProps {
  label: string;
  value: string | number | null;
  options: { value: string | number; label: string }[];
  onChange: (event: SelectChangeEvent<string | number>) => void;
  required?: boolean;
}

const SelectField: React.FC<SelectFieldProps> = ({ label, value, options, onChange, required }) => {
  const theme = useTheme();
  return (
    <FormControl fullWidth sx={{ marginBottom: 2 }}>
      {/* Подсказка (InputLabel) */}
      <InputLabel
        sx={{
          color: theme.palette.text.primary,
          backgroundColor: 'transparent', 
          paddingX: '4px', 
          '&.Mui-focused': {
            color: theme.palette.text.primary,
            backgroundColor: 'white', 
            transform: 'translate(14px, -6px) scale(0.75)', 
          },
          '&.MuiFormLabel-filled': {
            color: theme.palette.text.primary,
            backgroundColor: 'white', 
            transform: 'translate(14px, -6px) scale(0.75)', 
          },
          '&.MuiInputLabel-root': {
            transition: 'all 0.2s ease', 
          },
        }}
      >
        {label}
      </InputLabel>
      <Select
        value={value || ''}
        onChange={onChange}
        required={required}
        sx={{
          '& .MuiInputBase-input': {
            color: theme.palette.text.primary,
          },
          '& .MuiInputBase-root': {
            borderColor: theme.palette.text.primary,
            '&:hover': {
              borderColor: theme.palette.text.primary,
            },
          },
        }}
      >
        {options.map((option, index) => (
          <MenuItem key={index} value={option.value}>
            {option.label}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

export default SelectField;
