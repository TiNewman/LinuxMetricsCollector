import Link from 'next/link'
import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
//import Chart from "../components/CPUChart"
import dynamic from 'next/dynamic'

const Chart = dynamic(() => import('../components/CPUChart'), { ssr: false });

let socket

const cpuView = props => {
/*
  //use this to store the process list stuff
  const [process_list, setProcessList] = useState([])

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "process_list"}))
    };

    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)// might need to be e.data
      console.log(processJSON.process_list)
      setProcessList(processJSON.process_list)
      //setProcessList(e.data.process_list)
    }

    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }
  */
  const data = [30, 40, 45, 50, 49, 60, 70, 91]
  //const data = [{Usage:7.7905493}, {Usage:4.123711}, {Usage:5.1546392}, {Usage:5.050505}, {Usage:12.244898}, {Usage:4.1666665}, {Usage:4.0816326}, {Usage:13.131313}, {Usage:13.265306}]
  //<Chart data={data}/>
  //<div className="radial-progress text-primary ml-8" style={{"--value":data[0].Usage}}>{data[0].Usage}%</div>
  return (
    <div>
      <h1 className={processListStyles.h1}> CPU </h1>
      <Chart />
    </div>
  )
}

export default cpuView;