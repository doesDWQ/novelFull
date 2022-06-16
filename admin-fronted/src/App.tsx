import React from 'react';
import './App.css';
import {BrowserRouter, Route, Routes} from "react-router-dom";
import {HomePage} from "./components/homepage/HomePage";
import Users from "./components/homepage/Users/Users";
import AdminUsers from "./components/homepage/AdminUsers/AdminUsers";
import {Login} from "./components/login/Login";
import {LoginOut} from "./components/loginout/LoginOut";


function App() {
  return (
      <div>
          <React.StrictMode>
              <BrowserRouter>
                  <Routes>
                      <Route path="/" element={<HomePage />} >
                          <Route path="users" element={<Users />} />
                          <Route path="adminUsers" element={<AdminUsers />} />
                          <Route
                              index
                              element={<Users />}
                          />
                      </Route>
                      <Route path="login" element={<Login />} />
                      <Route path="loginOut" element={<LoginOut />} />
                      <Route
                          path="*"
                          element={
                              <main style={{ padding: "1rem" }}>
                                  <p>There's nothing here!</p>
                              </main>
                          }
                      />
                  </Routes>
              </BrowserRouter>
          </React.StrictMode>
      </div>
  );
}

export default App;
