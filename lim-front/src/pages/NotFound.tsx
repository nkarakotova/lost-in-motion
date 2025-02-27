import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Box, Typography } from '@mui/material';

const NotFound: React.FC = () => {
  const navigate = useNavigate();

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        backgroundColor: 'background.default',
        padding: 3,
      }}
    >

      <Typography variant="h1" sx={{ marginBottom: 3 }}>
        404
      </Typography>

      <Typography variant="h5" sx={{ marginBottom: 2 }}>
        Страница не найдена
      </Typography>

      <Button
        variant="contained"
        color="primary"
        onClick={() => navigate('/schedule')}
      >
        Вернуться на главную
      </Button>
    </Box>
  );
};

export default NotFound;
