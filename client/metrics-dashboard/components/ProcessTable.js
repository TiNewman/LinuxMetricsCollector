import TableScrollbar from 'react-table-scrollbar'
import Link from 'next/link'
import ProcessStyle from '../styles/Process_List.module.css'

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

const TableHeadItem = ({ item }) => <th className="bg-neutral text-base-100">{item.heading}</th>
const TableRow = ({ item, column }) => (
  <tr>
    {column.map((columnItem, index) => {
      return <td className="bg-primary text-neutral" key={item.PID + columnItem.value}>{item[`${columnItem.value}`]}</td>
    })}
  </tr>
)

export default ProcessTable;