import React, { useState } from 'react';
import Chart from 'react-apexcharts';

const CPUChart = () => {
  const [options, setOptions] = useState({
    chart: {
      id: 'line-chart',
    },
    xaxis: {
      categories: [
        '8:23:00',
        '8:23:30',
        '8:24:00',
        '8:24:30',
        '8:25:00',
        '8:25:30',
        '8:26:00',
      ],
    },
  });
  const [series, setSeries] = useState([
    {
      name: 'CPU Percentage Used',
      data: [7.7905493, 4.123711, 5.1546392, 5.050505, 12.244898, 4.1666665, 4.0816326],
    },
  ]);
  return (
    <div className='chart w-3/6 pt-7 mx-auto'>
      <Chart options={options} series={series} type='line'/>
    </div>
  );
}

export default CPUChart