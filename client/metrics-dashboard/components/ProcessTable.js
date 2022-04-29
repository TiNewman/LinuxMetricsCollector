import TableScrollbar from 'react-table-scrollbar'
import Link from 'next/link'
import ProcessStyle from '../styles/Process_List.module.css'

/**
 * ProcessTable
 *
 * Table that is used to display the process list data. Maps over the column and data arrays passed in so app doesn't
 * crash before the data is received
 *
 * Params: data: array of JSON with data values, column: the names of the column
 *
 * Return: html table for displaying process data
 */
const ProcessTable = ({ data, column }) => (
     <div className="overflow-x-auto flex flex-col p-5 justify-center items-center">
       <div className="block p-2 shadow-lg shadow-primary z-50">
       <TableScrollbar height="70vh">
           <table className="table table-zebra border-4 border-base-100 p-8">
             <thead className="border-4 border-base-100">
               <tr>
                 {column.map((item, index) => <TableHeadItem key={item.heading} item={item}/>)}
               </tr>
             </thead>
             <tbody>
               {data.map((item, index) => <TableRow key={item.PID} item={item} column={column} />)}
             </tbody>
           </table>
       </TableScrollbar>
       </div>
     </div>
)

/**
 * TableHeadItem
 *
 * Helper method that is called by the column map to set each column in the table with the appropriate name
 *
 * Param: item: JSON - contains the heading name to set
 *
 * Return: th element that is injected into table
 */
const TableHeadItem = ({ item }) => <th className="bg-neutral text-base-100">{item.heading}</th>

/**
 * TableRow
 *
 * Helper method that is called by the data map to set each row in the table with the appropriate data. Also maps over
 * the columns to put each item's data value in the correct column
 *
 * Param: item: JSON - contains the data to set, column - array of column names to map over to insert the data in
 *
 * Return: table column with data point properly placed under associated column
 */
const TableRow = ({ item, column }) => (
  <tr>
    {column.map((columnItem, index) => {
      return <td className="bg-primary text-neutral" key={item.PID + columnItem.value}>{item[`${columnItem.value}`]}</td>
    })}
  </tr>
)

export default ProcessTable;