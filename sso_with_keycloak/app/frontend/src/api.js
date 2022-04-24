import axios from "axios";
import { Issue } from "./issue.jsx";
import Keycloak from "keycloak-js";

const URL_PREFIX = "https://" + process.env.APP_DOMAIN_NAME + "/api";
const AUTHENTICATION_URL =
  "https://" + process.env.KEYCLOAK_DOMAIN_NAME + "/auth/";
const AUTHENTICATION_REALM = process.env.REALM;
const AUTHENTICATION_CLIENT_ID = process.env.CLIENT_ID;

let api;

export function createApi() {
  if (api == null) {
    if (process.env.REACT_FRONTEND_DUMMY_API == "true") {
      console.log("This application is standalone.");
      api = new DummyApi();
    } else {
      api = new Api();
    }
  }
  return api;
}

class Api {
  // --- 右記サイトから拝借 : https://stackoverflow.com/questions/60363513/how-to-execute-a-function-before-and-after-each-class-method-call //
  handler = {
    api: this,
    apply: function (target, scope, args) {
      const func_name = target.name;

      // 全ての API 通信の前にログイン状態を確認する。
      this.api.authenticate();
      const results = target.bind(this.api)(...args);

      return results;
    },
  };

  constructor(classAlias) {
    // Get all methods of choosen class
    let methods = Object.getOwnPropertyNames(Api.prototype);

    // Find and remove constructor as we don't need Proxy on it
    ["authenticate", "constructor"].forEach((methodName) => {
      let consIndex = methods.indexOf(methodName);
      if (consIndex > -1) methods.splice(consIndex, 1);
    });

    // Replace all methods with Proxy methods
    methods.forEach((methodName) => {
      this[methodName] = new Proxy(this[methodName], this.handler);
    });
    this.classAlias = classAlias;

    // Keycloak
    this.keycloak = new Keycloak({
      realm: AUTHENTICATION_REALM,
      url: AUTHENTICATION_URL,
      clientId: AUTHENTICATION_CLIENT_ID,
    });
  }
  // ---

  authenticate() {
    if (!this.authenticated) {
      this.keycloak
        .init({ onLoad: "login-required", checkLoginIframe: false })
        .then((authenticated) => {
          this.authenticated = authenticated;
        })
        .catch((e) => {
          console.log(e);
          alert("failed to initialize"); //TODO わかりやすいエラーメッセージにする
        });
    }
  }

  getUserRecords(setFunc) {
    //TODO
  }

  fight(userTypeValue, setRivalTypeValue, setIssueValue) {
    //TODO
  }
}

class DummyApi {
  constructor() {
    this.records = [];
  }

  getUserRecords(setRecords) {
    setRecords(this.records);
  }

  fight(userTypeValue, setRivalTypeValue, setIssueValue) {
    const rivalTypeValue = Math.floor(Math.random() * 3);
    let issueValue = null;
    if (userTypeValue == (rivalTypeValue + 1) % 3) {
      issueValue = Issue.Lose.value;
    } else if (userTypeValue == (rivalTypeValue + 2) % 3) {
      issueValue = Issue.Win.value;
    } else {
      issueValue = Issue.Draw.value;
    }
    setRivalTypeValue(rivalTypeValue);
    setIssueValue(issueValue);

    this.records.push({
      id: this.records.length + 1,
      issue: issueValue,
      type: userTypeValue,
      created_at: new Date(),
      updated_at: null,
      comment: null,
      is_edit: false,
    });
  }
}
