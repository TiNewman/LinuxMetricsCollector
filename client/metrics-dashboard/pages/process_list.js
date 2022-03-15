import Link from 'next/link'
import { useCallback, useEffect, useState } from 'react'
import io from 'Socket.IO-client'
import processListStyles from '../styles/Process_List.module.css'

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
  //const [message, setMessage] = useState('')
  //use this to store the process list stuff
  //const [process_list, setProcessList] = useState('')

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    await fetch('/api/socket');
    socket = io() //point socket to his here
    /*
    socket.on('connect', () => {
      setMessage('Connected')
    })

    socket.onmessage = (e) => {
      setMessage("Get message from server: " + e.data)
    };
    */
    socket.on('connect', () => {
      //as soon as connection happens, send request for process list
      socket.send(JSON.stringify({
            "request": "process_list"
      }))
    })

    socket.onmessage = (e) => {
      var processArray = JSON.parse(e)// might need to be e.data
      for (var i = 0; i < processArray.process_list.length; i++) {
          var process = processArray.process_list.counters[i];
          console.log(process.Name);
          createTableRow(process) //pass each process in and create a row for it
      }
    }

    return () => {
      socket.close()
    };
  }

  const createTableRow = (process) => {
    //declare the new row elements
    const tableRow = document.createElement("tr")
    const pid = document.createElement("td")
    const name = document.createElement("td")
    const cpu = document.createElement("td")
    const ram = document.createElement("td")
    const disk = document.createElement("td")
    const status = document.createElement("td")
    const uptime = document.createElement("td")

    //set the values of the process in the correct column
    pid.innerHTML = process.PID
    name.innerHTML = process.Name
    cpu.innerHTML = process.CPU
    ram.innerHTML = process.RAM
    disk.innerHTML = process.DISK
    status.innerHTML = process.STATUS
    uptime.innerHTML = process.UpTime

    //append the columns to the row
    tableRow.appendChild(pid)
    tableRow.appendChild(name)
    tableRow.appendChild(cpu)
    tableRow.appendChild(ram)
    tableRow.appendChild(disk)
    tableRow.appendChild(status)
    tableRow.appendChild(uptime)

    //append the row to the html table body
    const tableBody = html.getElementbyId("pTableBody")
    tableBody.appendChild(tableRow)
  }

  return html
}

export default processListView;