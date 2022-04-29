import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useEffect, useState } from 'react'
import io from 'socket.io-client'
import styles from '../styles/Home.module.css'
import Table from "../components/ProcessTable-Dashboard";
import Layout from "../components/Layout";

let socket

/**
 * Index
 *
 * Sends Websocket request when this page is loaded for all of the data. When request is received, it populates four
 * charts (Process List table, disk, CPU and RAM) in the HTML return statement and then displays the content on the
 * screen. Each chart contains data for all of the types of data collected (CPU, RAM, Processes, and Disk) and displays
 * a subset of each kind in the chart or table. Each table links to a page that contains that type of data in more detail
 *
 * Return: HTML that contains charts and tables for the different types of data metrics collected
 */
const Index = () => {

  //variables to set and store the incoming data
  const [process_list, setProcessList] = useState([])
  const [cpuData, setCPUData] = useState([])
  const [diskData, setDiskData] = useState([])
  const [ramData, setRAMData] = useState([])

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
    //make connection to websocket server and immediately send request for data
    socket.onopen = () => {
      console.log("Request being sent")
      socket.send(JSON.stringify({"request": "all"}))
    };

    //handles the incoming message from the websocket server
    socket.onmessage = (e) => {
      console.log("Received Message!: " + e.data)
      var processJSON = JSON.parse(e.data)
      setProcessList(processJSON.process_list)
      setCPUData([processJSON.cpu])
      setDiskData([processJSON.disk[0]])
      setRAMData([processJSON.memory])
    }

    //closes the websocket connection when the browser tab is closed
    return () => {
      console.log("closing socket")
      socket.send(JSON.stringify({"request": "stop"}))
      socket.close()
    };
  }

   /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS
   const response = {"process_list":[{"PID":1611,"Name":"systemd"},{"PID":1616,"Name":"(sd-pam)"},
                                     {"PID":1635,"Name":"gnome-keyring-d"},{"PID":1649,"Name":"gdm-wayland-ses"},
                                     {"PID":1652,"Name":"dbus-broker-lau"},{"PID":1654,"Name":"dbus-broker"},
                                     {"PID":1656,"Name":"gnome-session-b"}]}
   const process_list = response.process_list
   //console.log(process_list)



   const cpuData = [{Usage:37.7905493}]
   const diskData = [{"Name":"/dev/nvme0n1p3","MountPoint":"/","Usage":2.060362882143396,"Size":510405.902336}]
   const ramData = [{Usage:13.7905493}]
  */

  //column names for process list table -- only want to display PID and process name. Other details on process_list view
  const column = [
      { heading: 'PID', value: 'PID' },
      { heading: 'Name', value:'Name' },
    ]

  //Uses ProcessTable-Dashboard for process data and then Maps over cpuHData, ramData and diskData and returns html that
  //is injected into app view to avoid crashing before data is received
  return (
     <div>
       <div className="float-left mt-10 pt-10 pl-36">
         <h1 className={styles.h1}> Process List </h1>
         <Table data={process_list} column={column}/>
       </div>
       {cpuData.map((item, index) =>
       <Link href="/cpu">
         <div className="float-left mt-10 pt-10 pl-10">
           <h1 className={styles.h1}> CPU Usage </h1>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-neutral border-4 border-primary bg-primary hover:bg-base-100 hover:border-base-100"
                  style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
           </div>
         </div>
       </Link>
       )}
      {diskData.map((item, index) =>
       <Link href="/disk">
         <div className="float-left mt-10 pt-10 pl-16">
           <h2 className={styles.h2}>Root Disk: {item.Name}</h2>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-neutral border-4 border-primary bg-primary hover:bg-base-100 hover:border-base-100 "
                  style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
           </div>
         </div>
       </Link>
      )}
       {ramData.map((item, index) =>
       <Link href="/memory">
         <div className="float-left mt-10 pt-10 pl-16">
           <h1 className={styles.h1}> RAM Usage </h1>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-neutral border-4 border-primary bg-primary hover:bg-base-100 hover:border-base-100"
                  style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
           </div>
         </div>
       </Link>
      )}
     </div>
  )
}

export default Index;