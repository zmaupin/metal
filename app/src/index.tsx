import * as React from 'react';
import * as ReactDOM from 'react-dom';

const Main: () => JSX.Element = ():JSX.Element => (
  <h1>Hello World</h1>
);

ReactDOM.render(<Main />, document.getElementById('root'));
