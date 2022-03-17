import Link from 'next/link'
import { useCallback, useEffect, useState } from 'react'
import io from 'socket.io-client'
import processListStyles from '../styles/Process_List.module.css'
import Table from "../components/Table";

let socket
let html = (
 <div>
   <h1 className={processListStyles.h1}> Process List </h1>
   <div className="overflow-x-auto flex flex-col justify-center items-center">
     <table className="table">
       <thead>
         <tr>
           <th>PID</th>
           <th>Name</th>
           <th>CPU</th>
           <th>RAM</th>
           <th>Disk</th>
           <th>Status</th>
           <th>Execution Time</th>
         </tr>
       </thead>
       <tbody id="pTableBody">
       </tbody>
     </table>
   </div>
 </div>
)

const processListView = props => {
  //use this to store the process list stuff
  /*
  const [process_list, setProcessList] = useState('')

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "process_list"}))
    };

    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processArray = JSON.parse(e.data)// might need to be e.data
      for (var i = 0; i < processArray.process_list.length; i++) {
          var process = processArray.process_list[i];
          console.log(process.Name);
          createTableRow(process) //pass each process in and create a row for it
      }
    }

    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      setProcessList(e.data)
    }

    return () => {
      socket.close()
    };
  }*/

  const response = {"process_list":[{"PID":1611,"Name":"systemd","CPUUtilization":0.006125288,"RAMUtilization":14.983168,"DiskUtilization":37.363712,"Status":"S","ExecutionTime":6367.0474},{"PID":1616,"Name":"(sd-pam)","CPUUtilization":0,"RAMUtilization":6.770688,"DiskUtilization":0,"Status":"S","ExecutionTime":6367.0474},{"PID":1635,"Name":"gnome-keyring-d","CPUUtilization":0.0039270837,"RAMUtilization":7.708672,"DiskUtilization":0.106496,"Status":"S","ExecutionTime":6366.0474},{"PID":1649,"Name":"gdm-wayland-ses","CPUUtilization":0,"RAMUtilization":6.004736,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1652,"Name":"dbus-broker-lau","CPUUtilization":0.00015708334,"RAMUtilization":4.657152,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1654,"Name":"dbus-broker","CPUUtilization":0.00989625,"RAMUtilization":5.541888,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1656,"Name":"gnome-session-b","CPUUtilization":0.00015708334,"RAMUtilization":18.190336,"DiskUtilization":0.258048,"Status":"S","ExecutionTime":6366.0474}]}
  const process_list = [{"PID":1611,"Name":"systemd","CPUUtilization":0.006125288,"RAMUtilization":14.983168,"DiskUtilization":37.363712,"Status":"S","ExecutionTime":6367.0474},{"PID":1616,"Name":"(sd-pam)","CPUUtilization":0,"RAMUtilization":6.770688,"DiskUtilization":0,"Status":"S","ExecutionTime":6367.0474},{"PID":1635,"Name":"gnome-keyring-d","CPUUtilization":0.0039270837,"RAMUtilization":7.708672,"DiskUtilization":0.106496,"Status":"S","ExecutionTime":6366.0474},{"PID":1649,"Name":"gdm-wayland-ses","CPUUtilization":0,"RAMUtilization":6.004736,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1652,"Name":"dbus-broker-lau","CPUUtilization":0.00015708334,"RAMUtilization":4.657152,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1654,"Name":"dbus-broker","CPUUtilization":0.00989625,"RAMUtilization":5.541888,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1656,"Name":"gnome-session-b","CPUUtilization":0.00015708334,"RAMUtilization":18.190336,"DiskUtilization":0.258048,"Status":"S","ExecutionTime":6366.0474}]
  console.log(process_list)

  const column = [
    { heading: 'PID', value: 'PID' },
    { heading: 'Name', value:'Name' },
    { heading: 'CPU Utilization', value:'CPUUtilization'},
    { heading: 'RAM Utilization', value:'RAMUtilization'},
    { heading: 'Disk Utilization', value:'DiskUtilization'},
    { heading: 'Status', value:'Status'},
    { heading: 'Up Time', value:'ExecutionTime' },
  ]

  return (
    <Table data={process_list} column={column}/>
  )
}

export default processListView;