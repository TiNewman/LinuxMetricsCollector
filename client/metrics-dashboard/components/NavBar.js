import Link from 'next/link'

const NavBar = () => (
<div class="navbar bg-base-200 rounded-box">
  <div class="flex-1 px-2 lg:flex-none">
    <div class="flex items-stretch">
          <div class="dropdown">
            <label tabindex="0" class="btn btn-ghost rounded-btn">Linux Metrics Collector</label>
            <ul tabindex="0" class="menu dropdown-content p-2 shadow bg-secondary rounded-box w-52 mt-4">
              <li><Link href="/">Dashboard</Link></li>
              <li><Link href="/process_list">Process List</Link></li>
            </ul>
          </div>
        </div>
  </div>
  <div class="flex justify-end flex-1 px-2">
  </div>
</div>
);

export default NavBar;