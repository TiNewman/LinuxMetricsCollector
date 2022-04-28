import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
//import Table from "../components/ProcessTable"

let socket

const diskView = props => {
  //use this to store the process list stuff
  const [diskUsage, setDiskUsage] = useState([])

  useEffect(() => {
      socket = new WebSocket("ws://localhost:8080/ws");
      socketInitializer()
      return () => {
        console.log("closing socket")
        socket.send(JSON.stringify({"request": "stop"}))
        socket.close()
      };
     }, [])

  const socketInitializer = async () => {
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
  //const diskData = [{Usage:67.7905493}]
  const diskUsage = [{"Name":"/dev/nvme0n1p3","MountPoint":"/","Usage":2.060362882143396,"Size":510405.902336},{"Name":"/dev/nvme0n1p3","MountPoint":"/home","Usage":2.060362882143396,"Size":510405.902336},{"Name":"/dev/nvme0n1p2","MountPoint":"/boot","Usage":25.681494412006668,"Size":1020.70272},{"Name":"/dev/nvme0n1p1","MountPoint":"/boot/efi","Usage":2.3118672372403726,"Size":627.900416}]
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
        {diskUsage.map((item, index) => <div className="overflow-x-center flex flex-col w-1/6 mx-auto ml-24 mt-5 float-left px-4">
          <h2 className={processListStyles.h2}>{item.Name}</h2>
          <div className="block overflow-x-center flex flex-cols mt-5 p-5 shadow-lg shadow-primary">
            <div className="radial-progress text-neutral border-4 border-primary self-center bg-primary " style={{"--value":item.Usage.toFixed(2), "--size":"10rem"}}>{item.Usage.toFixed(2)}%</div>
          </div>
          <p className="pl-8 pt-8">Mount Point: {item.MountPoint}</p>
          <p className="pl-8">Size: {item.Size.toFixed(2)} MB</p>
        </div>)}
      </div>
  )
}

export default diskView;