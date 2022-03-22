import Head from 'next/head'
import Image from 'next/image'
import { useEffect, useState } from 'react'
import io from 'socket.io-client'
import styles from '../styles/Home.module.css'
import Table from "../components/Table";
//This will be the file for the main dashboard view with all of the elements

// pages/index.js

import Layout from "../components/Layout";

let socket

const Index = () => {
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
       setProcessList(processJSON.process_list)
     }

     return () => {
       console.log("closing socket")
       socket.send(JSON.stringify({"request": "stop"}))
       socket.close()
     };
   }
  */
     const response = {"process_list":[{"PID":1611,"Name":"systemd"},{"PID":1616,"Name":"(sd-pam)"},{"PID":1635,"Name":"gnome-keyring-d"},{"PID":1649,"Name":"gdm-wayland-ses"},{"PID":1652,"Name":"dbus-broker-lau"},{"PID":1654,"Name":"dbus-broker"},{"PID":1656,"Name":"gnome-session-b"}]}
     const process_list = response.process_list
     console.log(process_list)


   const column = [
       { heading: 'PID', value: 'PID' },
       { heading: 'Name', value:'Name' },
     ]

  return (
    <Table data={process_list} column={column}/>
  )
}

export default Index;