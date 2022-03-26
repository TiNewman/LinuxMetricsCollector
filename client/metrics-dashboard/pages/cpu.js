import Link from 'next/link'
import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
import Table from "../components/ProcessTable";

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

  return (
    <div>
      <h1 className={processListStyles.h1}> CPU </h1>
    </div>
  )
}

export default cpuView;