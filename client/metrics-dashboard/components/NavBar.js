import Link from 'next/link'

const NavBar = () => (
<div className="navbar bg-neutral rounded-box shadow-md shadow-accent">
  <div className="flex-1 px-2 lg:flex-none">
    <div className="flex items-stretch">
          <div className="dropdown">
            <label tabIndex="0" className="btn btn-ghost rounded-btn hover:bg-accent hover:text-neutral text-base-100">Linux Metrics Collector</label>
            <ul tabIndex="0" className="menu dropdown-content p-4 shadow-lg shadow-accent bg-neutral text-base-100 rounded-box w-52 mt-8">
              <li className="hover:bg-accent hover:text-neutral"><Link href="/">Dashboard</Link></li>
              <li className="hover:bg-accent hover:text-neutral"><Link href="/process_list">Process List</Link></li>
            </ul>
          </div>
        </div>
  </div>
  <div className="flex justify-end flex-1 px-2">
  </div>
</div>
);

export default NavBar;