// components/Layout.js

import Header from "./Header";
import NavBar from "./NavBar";
import Head from "next/head";

const Layout = ({children}) => (
  <div className="Layout">
    <Head>
      <title>Metrics Collector</title>
    </Head>
    <NavBar />
    {children}
  </div>
);

export default Layout;