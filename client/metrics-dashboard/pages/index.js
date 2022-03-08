import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'

//This will be the file for the main dashboard view with all of the elements

// pages/index.js

import Layout from "../components/Layout";

const Index = () => (
  <div>
    <h1 className={styles.h1}> Welcome to Linux Metrics Dashboard!</h1>
     <div className="overflow-x-auto flex flex-col justify-center items-center">
       <table class="table">
         <thead>
           <tr>
             <th>PID</th>
             <th>Name</th>
           </tr>
         </thead>
         <tbody>
           <tr className="hover">
             <th>1</th>
             <td></td>
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
     </div>
  </div>
);

export default Index;