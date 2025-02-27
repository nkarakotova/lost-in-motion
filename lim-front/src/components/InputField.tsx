import React from 'react';
import { TextField } from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface InputFieldProps {
  label: string;
  name: string;
  type?: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  required?: boolean;
  inputProps?: React.InputHTMLAttributes<HTMLInputElement>;
}


const InputField: React.FC<InputFieldProps> = ({
  label,
  name,
  type = 'text',
  value,
  onChange,
  required = false,
  inputProps,
}) => {
  const theme = useTheme();
  return (
    <TextField
      fullWidth
      label={label}
      name={name}
      type={type}
      variant="outlined"
      value={value}
      onChange={onChange}
      required={required}
      inputProps={inputProps}
      sx={{
        marginBottom: 2,
        '& .MuiInputLabel-root': {
          color: theme.palette.text.primary,
        },
        '& .MuiInputLabel-root.Mui-focused': {
          color: theme.palette.text.primary,
        },
        '& .MuiInputBase-input': {
          color: theme.palette.text.primary,
        },
      }}
    />
  );
};


export default InputField;
