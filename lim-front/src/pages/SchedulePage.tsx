import React from 'react';
import Header from '../components/Header';
import ScheduleTable from '../components/ScheduleTable';
import { ScheduleService } from '../api/services/scheduleService';
import { Box } from '@mui/material';

export const fetchScheduleData = async () => {
  const scheduleService = new ScheduleService();
  return scheduleService.getSchedule();
};

const SchedulePage: React.FC = () => {
  return (
    <Box display="flex" flexDirection="column" minHeight="100vh">
      <Header />
      <Box flexGrow={1} p={3}>
        <ScheduleTable fetchData={fetchScheduleData} />
      </Box>
    </Box>
  );
};

export default SchedulePage;
