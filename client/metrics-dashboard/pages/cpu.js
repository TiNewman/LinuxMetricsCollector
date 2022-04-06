import Link from 'next/link'
import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
import dynamic from 'next/dynamic'

const Chart = dynamic(() => import('react-apexcharts'), { ssr: false });

let socket
let dataArray = []
let categoriesArray = []

const cpuView = props => {
  //use this to store the CPU data
  const [options, setOptions] = useState({})
  const [series, setSeries] = useState([])

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "cpu"}))
    };

    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)// might need to be e.data
      console.log("Usage JSON:", processJSON.cpu.Usage)
      if(dataArray.length == 10){
        dataArray.shift()
        categoriesArray.shift()
      }
      dataArray.push(processJSON.cpu.Usage)

      //get current time
      var timestampInMilliseconds = Date.now();
      var timestampInSeconds = Date.now() / 1000; // A float value not an integer.
          timestampInSeconds = Math.floor(Date.now() / 1000); // Floor it to get the seconds.
      var time = new Date(timestampInSeconds).toISOString().substr(11, 8);
      categoriesArray.push(time)

      setOptions({
        chart: {
          id: 'line-chart',
        },
        xaxis: {
          categories: categoriesArray,
        },
      })
      setSeries([{
        name: 'CPU Percentage Used',
          data: dataArray,
        },
      ])

    }

    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }

  /*const [options, setOptions] = useState({
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
          '8:26:30',
          '8:27:00',
          '8:27:30',
        ],
      },
    });

    const [series, setSeries] = useState([
        {
          name: 'CPU Percentage Used',
          data: [7.7905493, 4.123711, 5.1546392, 5.050505, 12.244898, 4.1666665, 4.0816326, 7.7905493, 4.123711, 5.1546392],
        },
    ]); */


  //const incoming = [{Usage:7.7905493}, {Usage:4.123711}, {Usage:5.1546392}, {Usage:5.050505}, {Usage:12.244898}, {Usage:4.1666665}, {Usage:4.0816326}, {Usage:13.131313}, {Usage:13.265306}, {Usage:5.1546392}, {Usage:4.123711}]
  //var dataArray = []
  //var categoriesArray = []

  /*for (i = 0; i < 11; i++) {
    if(dataArray.length == 10){
      dataArray.shift()
      categoriesArray.shift()
    }
    incoming.shift() //won't need this
    //console.log(incoming[0].Usage)
    dataArray.push(incoming[0].Usage)

    //get current time
    var timestampInMilliseconds = Date.now();
    var timestampInSeconds = Date.now() / 1000; // A float value; not an integer.
        timestampInSeconds = Math.floor(Date.now() / 1000); // Floor it to get the seconds.
        timestampInSeconds = Date.now() / 1000 | 0; // Also you can do floor it like this.
        //timestampInSeconds = Math.round(Date.now() / 1000);
    var time = new Date(timestampInSeconds).toISOString().substr(11, 8);
    categoriesArray.push(time)

    setOptions({
                 chart: {
                   id: 'line-chart',
                 },
                 xaxis: {
                   categories: categoriesArray,
                 },
               })

    setSeries([{
                 name: 'CPU Percentage Used',
                 data: dataArray,
               },
              ])
    console.log(i)
  }*/

  //<div className="radial-progress text-primary ml-8" style={{"--value":data[0].Usage}}>{data[0].Usage}%</div>
  return (
    <div>
      <h1 className={processListStyles.h1}> CPU </h1>
      <div className='chart w-3/6 pt-7 mx-auto'>
            <Chart options={options} series={series} type='line'/>
      </div>
    </div>
  )
}

export default cpuView;