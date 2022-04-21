function containsAll(string, array) {
  for (var i=0; i < array.length; i++) {
    if (string.indexOf( array[i] ) == -1 ) {
      return false;
    }
  }
  return true;
}