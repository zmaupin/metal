import * as React from 'react';
import * as ReactDOM from 'react-dom';
import './bootstrap.tsx';
import CssBaseline from '@material-ui/core/CssBaseline';
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles';

import MAppBar from './components/MAppBar';


const theme = createMuiTheme({
  palette: {
    type: 'dark',
  },
});

const Main = ():JSX.Element => (
  <div id='root'>
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <MAppBar />
    </MuiThemeProvider>
  </div>
);

ReactDOM.render(<Main />, document.getElementById('root'));
