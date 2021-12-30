const Issue = {
  Win: { label: "勝ち", value: 0 },
  Lose: { label: "負け", value: 1 },
  Draw: { label: "引き分け", value: 2 },
};

function getIssueByValue(value) {
  switch (value) {
    case Issue.Win.value:
      return Issue.Win;
    case Issue.Lose.value:
      return Issue.Lose;
    case Issue.Draw.value:
      return Issue.Draw;
    default:
      return null;
  }
}

export { Issue, getIssueByValue };
