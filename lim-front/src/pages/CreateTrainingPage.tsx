import React, { useState, useEffect } from 'react';
import { HallService } from '../api/services/hallService';
import { CoachService } from '../api/services/coachesService';
import { Button, CircularProgress, Box, Typography, SelectChangeEvent } from '@mui/material';
import { TrainingInfo } from '../api/client/api';
import { TrainingService } from '../api/services/trainingsService';
import { useNavigate } from 'react-router-dom';
import InputField from '../components/InputField';
import SelectField from '../components/SelectField';
import Notification from '../components/Notification';

const CreateTrainingPage: React.FC = () => {
  const [date, setDate] = useState<string>('');
  const [groupedHalls, setGroupedHalls] = useState<any[]>([]);
  const [coaches, setCoaches] = useState<any[]>([]);
  const [selectedCoach, setSelectedCoach] = useState<string | null>(null);
  const [selectedHall, setSelectedHall] = useState<number | null>(null);
  const [selectedTime, setSelectedTime] = useState<string>('');
  const [placesNum, setPlacesNum] = useState<number>(15);
  const [trainingName, setTrainingName] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [openSnackbar, setOpenSnackbar] = useState<boolean>(false);
  const navigate = useNavigate();

  const coachService = new CoachService();
  const hallService = new HallService();

  useEffect(() => {
    const fetchCoaches = async () => {
      try {
        const fetchedCoaches = await coachService.getCoaches(setError);
        setCoaches(fetchedCoaches);
      } catch (err: any) {
        setError(err);
      }
    };

    fetchCoaches();
  }, []);

  useEffect(() => {
    if (date) {
      const fetchHalls = async () => {
        setLoading(true);
        setError(null);
        try {
          const fetchedHalls = await hallService.getHalls(date, setError);
          setGroupedHalls(fetchedHalls);
        } catch (err: any) {
          setGroupedHalls([]);
          setError(err);
        } finally {
          setLoading(false);
        }
      };

      fetchHalls();
    }
  }, [date]);

  
  const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setDate(e.target.value);
    setSelectedTime('');  
    setSelectedHall(null);  
  };

  
  const handleTimeChange = (e: SelectChangeEvent<string | number>) => {
    setSelectedTime(e.target.value as string);
    setSelectedHall(null);  
  };

  const handleSubmit = async () => {
    if (!selectedCoach || !selectedHall || !selectedTime || !date || !placesNum || !trainingName) {
      setError('Пожалуйста, заполните все поля');
      setOpenSnackbar(true);
      return;
    }

    const formattedDateTime = `${date}T${selectedTime}:00Z`;

    const trainingInfo: TrainingInfo = {
      coach_name: selectedCoach,
      hall_number: selectedHall,
      name: trainingName,
      date_time: formattedDateTime,
      places_num: placesNum,
    };

    try {
      const trainingService = new TrainingService();
      await trainingService.createTraining(trainingInfo, setError);

      setMessage('Тренировка успешно создана!');
      setError(null);
      setOpenSnackbar(true);
    } catch (err: any) {
      setMessage('');
      setError(err);
      setOpenSnackbar(true);
    }
  };

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        padding: 3,
        backgroundColor: 'background.default',
      }}
    >
      <Button
        variant="text"
        color="inherit"
        onClick={() => navigate('/schedule')}
        sx={{
          position: 'absolute',
          top: 16,
          left: 16,
        }}
      >
        Домой
      </Button>
      <Typography variant="h4" gutterBottom sx={{ textAlign: 'center' }}>
        Создание тренировки
      </Typography>

      <form
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit();
        }}
        style={{ width: '100%', maxWidth: 400 }}
      >
        <InputField
          label="Название тренировки"
          name="trainingName"
          value={trainingName}
          onChange={(e) => setTrainingName(e.target.value)}
          required
        />

        <InputField
          label=""
          name="date"
          type="date"
          value={date}
          onChange={handleDateChange} 
          required
        />

        {loading && <CircularProgress size={24} color="secondary" sx={{ marginBottom: 2 }} />}

        <SelectField
          label="Выберите тренера"
          value={selectedCoach}
          options={coaches.map((coach) => ({ value: coach.name, label: coach.name }))}
          onChange={(e: SelectChangeEvent<string | number>) => setSelectedCoach(e.target.value as string)}
          required
        />

        <SelectField
          label="Выберите время"
          value={selectedTime}
          options={groupedHalls
            .map((group) => ({ value: group.time, label: group.time })) 
            .filter((option, index, self) => self.findIndex(o => o.value === option.value) === index)} 
          onChange={handleTimeChange} 
          required
        />

        <SelectField
          label="Выберите зал"
          value={selectedHall}
          options={groupedHalls
            .filter(group => group.time === selectedTime) 
            .flatMap(group => group.halls) 
            .map((hall) => ({ value: hall.number, label: `Зал ${hall.number}` }))}
          onChange={(e: SelectChangeEvent<string | number>) => setSelectedHall(Number(e.target.value))}
          required
        />

        <InputField
          label="Количество мест"
          name="placesNum"
          type="number"
          value={placesNum.toString()}
          onChange={(e) => setPlacesNum(Number(e.target.value))}
          inputProps={{ min: 1 }}
          required
        />

        <Button
          variant="contained"
          color="primary"
          onClick={handleSubmit}
          disabled={loading}
          sx={{ width: '100%', marginTop: 2 }}
        >
          {loading ? 'Загрузка...' : 'Создать тренировку'}
        </Button>
      </form>

      <Notification
        open={openSnackbar}
        message={error || message}
        severity={error ? 'error' : 'success'}
        onClose={() => setOpenSnackbar(false)}
      />
    </Box>
  );
};

export default CreateTrainingPage;
