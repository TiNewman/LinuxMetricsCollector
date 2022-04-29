import Link from 'next/link'
import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
import dynamic from 'next/dynamic'

//dynamic import for Apex charts
const Chart = dynamic(() => import('react-apexcharts'), { ssr: false });

//websocket variables
let socket
let dataArray = []
let categoriesArray = []

/**
 * cpuView
 *
 * Sends Websocket request when this page is loaded for CPU data. When request is received, it populates the apex chart in
 * the HTML return statement and then displays the content on the screen.
 *
 * Return: HTML that contains a label for the CPU page and the CPU chart
 */
const cpuView = props => {
  //variables to set and store the incoming CPU data
  const [options, setOptions] = useState({})
  const [series, setSeries] = useState([])

  //initialize Websocket and close the connection when this page is unmounted from the view
  useEffect(() => {
    socket = new WebSocket("ws://localhost:8080/ws");
    socketInitializer()
    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }, [])

  //initializes methods for the websocket
  const socketInitializer = async () => {
    //make connection to websocket server and immediately send request for CPU data
    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "cpu"}))
    };

    //handles the incoming message from the websocket server
    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)
      //only want to display 10 data points in chart, so if length is already 10, then shift the array every time a new
      //point comes in
      if(dataArray.length == 10){
        dataArray.shift()
        categoriesArray.shift()
      }
      //push the new data point onto the array
      dataArray.push(processJSON.cpu.Usage.toFixed(2))

      //get current time and push it to the categories array for the x-axis
      var curtime = new Date()
      var time = "" + curtime.getHours() + ":" + curtime.getMinutes() + ":" + curtime.getSeconds()
      categoriesArray.push(time)

      //Sets the axes and other styling for the line chart
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

      //passes the data points into the chart
      setSeries([
        {
          name: 'CPU Percentage Used',
          data: dataArray,
        },
      ])
    }
    //closes the websocket connection when the browser tab is closed
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

  //HTML that is injected into the app view
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