import Head from 'next/head'
import Image from 'next/image'
import { useEffect, useState } from 'react'
import io from 'Socket.IO-client'
import styles from '../styles/Home.module.css'
//This will be the file for the main dashboard view with all of the elements

// pages/index.js

import Link from 'next/link'
import Layout from "../components/Layout";

let socket

const Index = () => {
  const [message, setMessage] = useState('')
  const [process_list, setProcessList] = useState('')

  useEffect(() => socketInitializer(), [])

  const socketInitializer = async () => {
    await fetch('/api/socket');
    socket = io() //point socket to his here

    //might need to use the set message to store the incoming data -- declare variable?
    socket.on('connect', () => {
      setMessage('Connected')
    })

    socket.onmessage = (e) => {
          setMessage("Get message from server: " + e.data)
    };

    socket.on('connect', () => {
          setProcessList()
        })

    socket.onmessage = (e) => {
          setProcessList("Get message from server: " + e.data)
    };

  }

  return (
    <div>
        <h1 className={styles.h1}> Welcome to Linux Metrics Dashboard!</h1>
         <div className="overflow-x-auto flex flex-col justify-center items-center">
           <Link href="/process_list">
           <table className="table">
             <thead>
               <tr>
                 <th>PID</th>
                 <th>Name</th>
               </tr>
             </thead>
             <tbody>
               <tr className="hover">
                 <th>1</th>
                 <td>{message}</td>
               </tr>
               <tr className="hover">
                 <th>2</th>
                 <td></td>
               </tr>
               <tr className="hover">
                 <th>3</th>
                 <td></td>
               </tr>
               <tr className="hover">
                 <th>4</th>
                 <td></td>
               </tr>
               <tr className="hover">
                 <th>5</th>
                 <td></td>
               </tr>
             </tbody>
           </table>
           </Link>
         </div>
      </div>
  )
}

export default Index;