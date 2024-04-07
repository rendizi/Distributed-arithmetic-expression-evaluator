import React from "react";
import Default from "./Main/Default";
import Operations from "./Main/Operations";
import Expressions from "./Main/Expressions";
import Agents from "./Main/Agents";
import Expression from "./Main/Expression";

function Main(props){
    return (
        <main>
            <div id="app">
                <div id="star-container">
                    <div id="star-pattern"></div>
                    <div id="star-gradient-overlay"></div>
                    <div id="center-box"></div>
                    {props.window === "operations" ? <Operations /> : props.window === "agents" ? <Agents />
                     : props.window === "expressions" ? <Expressions /> : props.window === "expression" ?
                      <Expression expression={props.expression} result = {props.result} time = {props.time}/> : <Default /> }
                </div>
            </div>
        </main>
    );
}

export default Main 