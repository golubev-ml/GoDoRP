import React from "react";
import ReactDOM from "react-dom";
import { Router, Route, browserHistory } from "react-router";
import App from "./components/App";
import Counters from "./components/counters";

ReactDOM.render(
  /*<Counters />,*/
  <Router history={browserHistory}>
    <Route path="/" component={App} />
  </Router>,
  document.getElementById("root")
);
