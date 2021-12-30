const Type = {
  Rock: { label: "ぐー", value: 0 },
  Scissors: { label: "ちょき", value: 1 },
  Paper: { label: "ぱー", value: 2 },
};

const TypeArray = Object.entries(Type).map((type) => {
  return type[1];
});

function getTypeByValue(value) {
  switch (value) {
    case Type.Rock.value:
      return Type.Rock;
    case Type.Scissors.value:
      return Type.Scissors;
    case Type.Paper.value:
      return Type.Paper;
    default:
      return null;
  }
}

export { Type, TypeArray, getTypeByValue };
