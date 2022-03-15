import Link from 'next/link'

const NavBar = () => (
<div className="navbar bg-base-200 rounded-box">
  <div className="flex-1 px-2 lg:flex-none">
    <div className="flex items-stretch">
          <div className="dropdown">
            <label tabIndex="0" className="btn btn-ghost rounded-btn">Linux Metrics Collector</label>
            <ul tabIndex="0" className="menu dropdown-content p-2 shadow bg-secondary rounded-box w-52 mt-4">
              <li><Link href="/">Dashboard</Link></li>
              <li><Link href="/process_list">Process List</Link></li>
            </ul>
          </div>
        </div>
  </div>
  <div className="flex justify-end flex-1 px-2">
  </div>
</div>
);

export default NavBar;