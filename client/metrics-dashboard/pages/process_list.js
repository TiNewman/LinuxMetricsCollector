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
       <tbody className="scroll-smooth" id="pTableBody">
       </tbody>
     </table>
   </div>
 </div>
)

const processListView = props => {
  //use this to store the process list stuff
  const [process_list, setProcessList] = useState('')

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      socket.send(JSON.stringify({"request": "process_list"}))
    };

    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)// might need to be e.data
      setProcessList(processJSON.process_list)
    }

    return () => {
      console.log("closing socket")
      socket.close()
    };
  }

  /*
  const response = {"process_list":[{"PID":1611,"Name":"systemd","CPUUtilization":0.006125288,"RAMUtilization":14.983168,"DiskUtilization":37.363712,"Status":"S","ExecutionTime":6367.0474},{"PID":1616,"Name":"(sd-pam)","CPUUtilization":0,"RAMUtilization":6.770688,"DiskUtilization":0,"Status":"S","ExecutionTime":6367.0474},{"PID":1635,"Name":"gnome-keyring-d","CPUUtilization":0.0039270837,"RAMUtilization":7.708672,"DiskUtilization":0.106496,"Status":"S","ExecutionTime":6366.0474},{"PID":1649,"Name":"gdm-wayland-ses","CPUUtilization":0,"RAMUtilization":6.004736,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1652,"Name":"dbus-broker-lau","CPUUtilization":0.00015708334,"RAMUtilization":4.657152,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1654,"Name":"dbus-broker","CPUUtilization":0.00989625,"RAMUtilization":5.541888,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1656,"Name":"gnome-session-b","CPUUtilization":0.00015708334,"RAMUtilization":18.190336,"DiskUtilization":0.258048,"Status":"S","ExecutionTime":6366.0474},{"PID":1700,"Name":"gnome-session-c","CPUUtilization":0,"RAMUtilization":5.095424,"DiskUtilization":0.008192,"Status":"S","ExecutionTime":6366.0474},{"PID":1701,"Name":"uresourced","CPUUtilization":0,"RAMUtilization":4.980736,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1704,"Name":"gnome-session-b","CPUUtilization":0.0017279168,"RAMUtilization":19.771393,"DiskUtilization":0.135168,"Status":"S","ExecutionTime":6366.0474},{"PID":1737,"Name":"gnome-shell","CPUUtilization":5.789935,"RAMUtilization":338.006,"DiskUtilization":7.02464,"Status":"S","ExecutionTime":6366.0474},{"PID":1873,"Name":"gvfsd","CPUUtilization":0.00062833336,"RAMUtilization":8.425472,"DiskUtilization":0.106496,"Status":"S","ExecutionTime":6366.0474},{"PID":1882,"Name":"gvfsd-fuse","CPUUtilization":0,"RAMUtilization":6.598656,"DiskUtilization":0.135168,"Status":"S","ExecutionTime":6366.0474},{"PID":1889,"Name":"at-spi-bus-laun","CPUUtilization":0,"RAMUtilization":7.745536,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1894,"Name":"dbus-broker-lau","CPUUtilization":0,"RAMUtilization":4.386816,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1895,"Name":"dbus-broker","CPUUtilization":0.00031416668,"RAMUtilization":2.92864,"DiskUtilization":0,"Status":"S","ExecutionTime":6366.0474},{"PID":1908,"Name":"xdg-permission-","CPUUtilization":0,"RAMUtilization":5.2224,"DiskUtilization":0,"Status":"S","ExecutionTime":6365.0474},{"PID":1911,"Name":"gnome-shell-cal","CPUUtilization":0.00047132405,"RAMUtilization":22.831104,"DiskUtilization":5.640192,"Status":"S","ExecutionTime":6365.0474},{"PID":1935,"Name":"evolution-sourc","CPUUtilization":0.0006284321,"RAMUtilization":28.033024,"DiskUtilization":0.978944,"Status":"S","ExecutionTime":6365.0474},{"PID":1939,"Name":"pipewire","CPUUtilization":2.1418536,"RAMUtilization":21.233664,"DiskUtilization":0.04096,"Status":"S","ExecutionTime":6365.0474}]}
  const responseArray = response.process_list
  console.log(responseArray)
  */

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