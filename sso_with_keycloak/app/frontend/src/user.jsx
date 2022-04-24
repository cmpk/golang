import React, { useEffect, useState } from "react";
import { getTypeByValue, TypeArray } from "./type.jsx";
import { Issue, getIssueByValue } from "./issue.jsx";
import { createApi } from "./api.js";

const TypeItem = (props) => {
  const backgroundColor =
    props.type.value === props.userTypeValue ? "#f88" : "#3cf";
  const fontColor = props.type.value === props.userTypeValue ? "#fff" : "#000";
  return (
    <div style={{ marginBottom: "5px" }}>
      <button
        style={{
          fontSize: "large",
          padding: "5px 0px",
          border: "none",
          borderRadius: "20px",
          backgroundColor: backgroundColor,
          color: fontColor,
          width: "6em",
        }}
        onClick={() => props.setUserTypeValue(props.type.value)}
      >
        {props.type.label}
      </button>
    </div>
  );
};

const Battle = (props) => {
  return (
    <div
      style={{
        borderRadius: "10px",
        backgroundColor: "#ffb",
        minWidth: "300px",
        padding: "0 20px 10px 20px",
      }}
    >
      <h2>あなたの手</h2>
      <div
        style={{
          display: "flex",
          justifyContent: "space-around",
        }}
      >
        <div style={{ textAlign: "center" }}>
          {TypeArray.map((type) => {
            return (
              <TypeItem
                key={type.value}
                type={type}
                userTypeValue={props.userTypeValue}
                setUserTypeValue={props.setUserTypeValue}
              />
            );
          })}
        </div>
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          <button
            style={{
              display: "inline",
              fontSize: "x-large",
              borderRadius: "15px",
              border: "none",
              padding: "10px 20px",
              backgroundColor: "#fb9",
            }}
            onClick={() => props.checkResult()}
            disabled={props.userTypeValue == null}
          >
            勝負！
          </button>
        </div>
      </div>
    </div>
  );
};

const Result = (props) => {
  const message = (issueValue) => {
    switch (issueValue) {
      case Issue.Win.value:
        return "勝ちました！";
      case Issue.Lose.value:
        return "負けました...";
      case Issue.Draw.value:
        return "引き分けです";
      default:
        return "";
    }
  };
  return (
    <div
      style={{
        borderRadius: "10px",
        backgroundColor: "#bff",
        minWidth: "300px",
        padding: "0 20px",
      }}
    >
      <h2>相手の手</h2>
      <div style={{ textAlign: "center" }}>
        <div style={{ fontSize: "x-large", marginBottom: "15px" }}>
          {props.rivalTypeValue != null &&
            getTypeByValue(props.rivalTypeValue).label}
        </div>
        <div>{message(props.issueValue)}</div>
      </div>
    </div>
  );
};

const Records = (props) => {
  if (props.records.length === 0) {
    return <div style={{ marginLeft: "20px" }}>なし</div>;
  }

  const records = props.records.slice().reverse();
  return (
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>勝ち負け</th>
          <th>あなたの手</th>
          <th>勝負した日時</th>
        </tr>
      </thead>
      <tbody id="record">
        {records.map((record, index) => {
          return (
            <tr key={records.length - (index + 1)}>
              <td>{record.id}</td>
              <td>{getIssueByValue(record.issue).label}</td>
              <td>{getTypeByValue(record.type).label}</td>
              <td>
                {record.created_at.toLocaleDateString()}{" "}
                {record.created_at.toLocaleTimeString()}
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

export function UserPage() {
  const [userTypeValue, setUserTypeValue] = useState();
  const [rivalTypeValue, setRivalTypeValue] = useState();
  const [issueValue, setIssueValue] = useState();
  const [records, setRecords] = useState([]);

  const api = createApi();
  useEffect(() => {
    api.getUserRecords(setRecords);
  }, []);

  const checkResult = () => {
    api.fight(userTypeValue, setRivalTypeValue, setIssueValue);
    api.getUserRecords(setRecords);
  };

  return (
    <div>
      <h1>じゃんけん</h1>

      <div
        style={{
          display: "flex",
          justifyContent: "space-around",
          maxWidth: "800px",
          marginBottom: "10px",
        }}
      >
        <Battle
          userTypeValue={userTypeValue}
          setUserTypeValue={setUserTypeValue}
          checkResult={checkResult}
        />
        <Result issueValue={issueValue} rivalTypeValue={rivalTypeValue} />
      </div>

      <h2>勝敗履歴</h2>
      <Records records={records} setRecords={setRecords} />
    </div>
  );
}
