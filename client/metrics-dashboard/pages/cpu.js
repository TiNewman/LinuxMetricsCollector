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
        markers: {
          size: 0,
          colors: ['#0D5090'],
        },
        stroke: {
          colors: ['#BFAFFF'],
        },
        xaxis: {
          categories: categoriesArray,
          labels: {
            show: true,
            align: 'right',
            minWidth: 0,
            maxWidth: 160,
            style: {
              colors: ['#0D5090', '#FFA207', '#0D5090', '#FFA207', '#0D5090', '#FFA207', '#0D5090','#FFA207','#0D5090','#FFA207'],
              fontSize: '14px',
              fontFamily: 'Helvetica, Arial, sans-serif',
              fontWeight: 400,
            },
            offsetX: 0,
            offsetY: 0,
            rotate: 0,
            },
            title: {
              text: 'TimeStamp',
              rotate: -90,
              offsetX: 0,
              offsetY: 0,
              style: {
                color: 'white',
                fontSize: '14px',
                fontFamily: 'Helvetica, Arial, sans-serif',
                fontWeight: 600,
              },
            }
          },
          yaxis: {
            show: true,
            showAlways: true,
            showForNullSeries: true,
            seriesName: 'CPU Percentage Used',
            logBase: 10,
            tickAmount: 5,
            min: 0,
            max: 100,
            labels: {
              show: true,
              align: 'right',
              minWidth: 0,
              maxWidth: 160,
              style: {
                colors: ['#FFA207', '#0D5090'],
                fontSize: '14px',
                fontFamily: 'Helvetica, Arial, sans-serif',
                fontWeight: 400,
              },
              offsetX: 0,
              offsetY: 0,
              rotate: 0,
            },
            title: {
              text: 'CPU Percentage Used',
              rotate: -90,
              offsetX: 0,
              offsetY: 0,
              style: {
                color: 'white',
                fontSize: '14px',
                fontFamily: 'Helvetica, Arial, sans-serif',
                fontWeight: 600,
              },
            }
          },
        })
      setSeries([
        {
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

  /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS 
  const [options, setOptions] = useState({
      chart: {
        id: 'line-chart',
      },
      markers: {
          size: 0,
          colors: ['#0D5090'],
      },
     stroke: {
       colors: ['#BFAFFF'],
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
        labels: {
          show: true,
          align: 'right',
          minWidth: 0,
          maxWidth: 160,
          style: {
            colors: ['#0D5090', '#FFA207', '#0D5090', '#FFA207', '#0D5090', '#FFA207', '#0D5090','#FFA207','#0D5090','#FFA207'],
            fontSize: '14px',
            fontFamily: 'Helvetica, Arial, sans-serif',
            fontWeight: 400,
          },
          offsetX: 0,
          offsetY: 0,
          rotate: 0,
        },
        title: {
          text: 'TimeStamp',
          rotate: -90,
          offsetX: 0,
          offsetY: 0,
          style: {
            color: 'white',
            fontSize: '14px',
            fontFamily: 'Helvetica, Arial, sans-serif',
            fontWeight: 600,
          },
        }
      },
       yaxis: {
            show: true,
            showAlways: true,
            showForNullSeries: true,
            seriesName: 'CPU Percentage Used',
            logBase: 10,
            tickAmount: 5,
            min: 0,
            max: 100,
            labels: {
                show: true,
                align: 'right',
                minWidth: 0,
                maxWidth: 160,
                style: {
                    colors: ['#FFA207', '#0D5090'],
                    fontSize: '14px',
                    fontFamily: 'Helvetica, Arial, sans-serif',
                    fontWeight: 400,
                },
                offsetX: 0,
                offsetY: 0,
                rotate: 0,
            },
            title: {
                text: 'CPU Percentage Used',
                rotate: -90,
                offsetX: 0,
                offsetY: 0,
                style: {
                    color: 'white',
                    fontSize: '14px',
                    fontFamily: 'Helvetica, Arial, sans-serif',
                    fontWeight: 600,
                },
            }
        },
    });

    const [series, setSeries] = useState([
        {
          name: 'CPU Percentage Used',
          data: [7.7905493, 4.123711, 5.1546392, 5.050505, 12.244898, 4.1666665, 4.0816326, 7.7905493, 4.123711, 5.1546392],
        },
    ]);
 */

  return (
    <div>
      <h1 className={processListStyles.h1}> CPU </h1>
      <div className='chart w-3/6 pt-7 mx-auto block p-2 shadow-lg shadow-primary'>
            <Chart options={options} series={series} type='line'/>
      </div>
    </div>
  )
}

export default cpuView;