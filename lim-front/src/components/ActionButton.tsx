import React, { useState } from 'react';
import { Button } from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface ActionButtonProps {
  label: string;
  onClick: () => void;
}

const ActionButton: React.FC<ActionButtonProps> = ({ label, onClick }) => {
  const theme = useTheme();
  const [clicked, setClicked] = useState(false); 

  const handleClick = () => {
    setClicked(true); 
    onClick(); 

    setTimeout(() => {
      setClicked(false);
    }, 500); 
  };

  return (
    <Button
      variant="contained"
      color={clicked ? 'secondary' : 'primary'} 
      onClick={handleClick}
      sx={{
        minWidth: '200px',
        margin: '5px',
        backgroundColor: clicked ? theme.palette.secondary.main : theme.palette.primary.main, 
      }}
    >
      {label}
    </Button>
  );
};

export default ActionButton;
