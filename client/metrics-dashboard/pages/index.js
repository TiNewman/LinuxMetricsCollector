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
   //const [process_list, setProcessList] = useState([])
   //const [cpuData, setCPUData] = useState([])
   //const [diskData, setDiskData] = useState([])
   //const [ramData, setRAMData] = useState([])

   //useEffect(() => socketInitializer(), [])

   /*const socketInitializer = async () => {
     const socket = new WebSocket("ws://localhost:8080/ws");

     socket.onopen = () => {
       console.log("Request being sent")
       //socket.send(JSON.stringify({"request": "process_list"}))
       socket.send(JSON.stringify({"request": "all"}))
     };

     socket.onmessage = (e) => {
       console.log("Received Message!: " + e.data)
       var processJSON = JSON.parse(e.data)// might need to be e.data
       setProcessList(processJSON.process_list)
       setCPUData(processJSON.cpu)
       setDiskData(processJSON.disk)
       setRAMData(processJSON.mem)
     }

     return () => {
       console.log("closing socket")
       socket.send(JSON.stringify({"request": "stop"}))
       socket.close()
     };
   }*/

   /** THIS IS THE MANUAL TEST DATA FOR CLIENT WEBSOCKETS */
   const response = {"process_list":[{"PID":1611,"Name":"systemd"},{"PID":1616,"Name":"(sd-pam)"},{"PID":1635,"Name":"gnome-keyring-d"},{"PID":1649,"Name":"gdm-wayland-ses"},{"PID":1652,"Name":"dbus-broker-lau"},{"PID":1654,"Name":"dbus-broker"},{"PID":1656,"Name":"gnome-session-b"}]}
   const process_list = response.process_list
   //console.log(process_list)


   const column = [
       { heading: 'PID', value: 'PID' },
       { heading: 'Name', value:'Name' },
     ]

   const cpuData = [{Usage:37.7905493}]
   const diskData = [{Usage:67.7905493}]
   const ramData = [{Usage:13.7905493}]
  /**/

  return (
     <div>
       <div className="float-left mt-10 pt-10 pl-56">
         <h1 className={styles.h1}> Process List </h1>
         <Table data={process_list} column={column}/>
       </div>
       <Link href="/cpu">
         <div className="float-left mt-10 pt-10 pl-6">
           <h1 className={styles.h1}> CPU Usage </h1>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-primary hover:text-base-100" style={{"--value":cpuData[0].Usage, "--size":"12rem"}}>{cpuData[0].Usage}%</div>
           </div>
         </div>
       </Link>
       <Link href="/disk">
         <div className="float-left mt-10 pt-10 pl-10">
           <h1 className={styles.h1}> Disk Usage </h1>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-primary hover:text-base-100" style={{"--value":diskData[0].Usage, "--size":"12rem"}}>{diskData[0].Usage}%</div>
           </div>
         </div>
       </Link>
       <Link href="/memory">
         <div className="float-left mt-10 pt-10 pl-10">
           <h1 className={styles.h1}> RAM Usage </h1>
           <div className="block p-5 shadow-lg shadow-primary hover:bg-primary">
             <div className="radial-progress text-primary hover:text-base-100" style={{"--value":ramData[0].Usage, "--size":"12rem"}}>{ramData[0].Usage}%</div>
           </div>
         </div>
       </Link>
     </div>
  )
}

export default Index;