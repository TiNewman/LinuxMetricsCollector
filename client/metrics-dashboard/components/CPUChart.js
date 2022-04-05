import React, { useState } from 'react';
import Chart from 'react-apexcharts';

const CPUChart = (options, series) => {
  const [options, setOptions] = useState({});
  const [series, setSeries] = useState([]);

  setOptions = {options}
  setSeries = {series}

  return (
    <div className='chart w-3/6 pt-7 mx-auto'>
      <Chart options={options} series={series} type='line'/>
    </div>
  );
}

export default CPUChart