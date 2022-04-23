import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
//import Table from "../components/ProcessTable"

let socket

const diskView = props => {
  //use this to store the process list stuff
  const [diskUsage, setDiskUsage] = useState([])

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "disk"}))
      console.log("sent disk request")
    };

    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)// might need to be e.data
      console.log(processJSON.disk)
      setDiskUsage(processJSON.disk)
    }

    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }

  /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS 
  const diskData = [{Usage:67.7905493}]
  */

  //<div className="radial-progress text-primary self-center" style={{"--value":diskUsage[0].Usage, "--size":"20rem"}}>{diskUsage[0].Usage}%</div>
      /*
      <div className="overflow-x-center flex flex-col w-1/4 mx-auto">
        <div className="block overflow-x-center flex flex-cols mt-5 p-5 shadow-lg shadow-primary">
          {diskUsage.map((item, index) => <div className="radial-progress text-primary self-center" style={{"--value":item.Usage, "--size":"20rem"}}>{item.Usage}%</div>)}
        </div>
      </div>
      */
  return (
    <div>
      <h1 className={processListStyles.h1}> Disk Usage </h1>
      {diskUsage.map((item, index) => <div className="overflow-x-center flex flex-col w-1/4 mx-auto">
        <div className="block overflow-x-center flex flex-cols mt-5 p-5 shadow-lg shadow-primary">
          <div className="radial-progress text-primary self-center" style={{"--value":item.Usage, "--size":"20rem"}}>{item.Usage}%</div>
        </div>
      </div>)}
    </div>
  )
}

export default diskView;