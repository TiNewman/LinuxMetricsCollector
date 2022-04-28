import Head from 'next/head'
import Image from 'next/image'
import Link from 'next/link'
import { useEffect, useState } from 'react'
import io from 'socket.io-client'
import styles from '../styles/Home.module.css'
import Table from "../components/ProcessTable-Dashboard";
//This will be the file for the main dashboard view with all of the elements

// pages/index.js

import Layout from "../components/Layout";

let socket
const Index = () => {

   //use this to store the process list stuff
   const [process_list, setProcessList] = useState([])
   const [cpuData, setCPUData] = useState([])
   const [diskData, setDiskData] = useState([])
   const [ramData, setRAMData] = useState([])

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
       console.log("Request being sent")
       //socket.send(JSON.stringify({"request": "process_list"}))
       socket.send(JSON.stringify({"request": "all"}))
     };

     socket.onmessage = (e) => {
       console.log("Received Message!: " + e.data)
       var processJSON = JSON.parse(e.data)// might need to be e.data
       setProcessList(processJSON.process_list)
       setCPUData([processJSON.cpu])
       setDiskData([processJSON.disk[0]])
       setRAMData([processJSON.memory])
     }

     return () => {
       console.log("closing socket")
       socket.send(JSON.stringify({"request": "stop"}))
       socket.close()
     };
   }

   /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS
   const response = {"process_list":[{"PID":1611,"Name":"systemd"},{"PID":1616,"Name":"(sd-pam)"},{"PID":1635,"Name":"gnome-keyring-d"},{"PID":1649,"Name":"gdm-wayland-ses"},{"PID":1652,"Name":"dbus-broker-lau"},{"PID":1654,"Name":"dbus-broker"},{"PID":1656,"Name":"gnome-session-b"}]}
   const process_list = response.process_list
   //console.log(process_list)



   const cpuData = [{Usage:37.7905493}]
   const diskData = [{"Name":"/dev/nvme0n1p3","MountPoint":"/","Usage":2.060362882143396,"Size":510405.902336}]
   const ramData = [{Usage:13.7905493}]
  */
  const column = [
      { heading: 'PID', value: 'PID' },
      { heading: 'Name', value:'Name' },
    ]

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
             <div className="radial-progress text-neutral border-4 border-primary bg-primary hover:bg-base-100 hover:border-base-100" style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
           </div>
         </div>
       </Link>
       )}
      {diskData.map((item, index) =>
       <Link href="/disk">
         <div className="float-left mt-10 pt-10 pl-16">
           <h2 className={styles.h2}>Root Disk: {item.Name}</h2>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-neutral border-4 border-primary bg-primary hover:bg-base-100 hover:border-base-100 " style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
           </div>
         </div>
       </Link>
      )}
       {ramData.map((item, index) =>
       <Link href="/memory">
         <div className="float-left mt-10 pt-10 pl-16">
           <h1 className={styles.h1}> RAM Usage </h1>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-neutral border-4 border-primary bg-primary hover:bg-base-100 hover:border-base-100" style={{"--value":item.Usage.toFixed(2), "--size":"12rem"}}>{item.Usage.toFixed(2)}%</div>
           </div>
         </div>
       </Link>
      )}
     </div>
  )
}

export default Index;