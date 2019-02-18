import * as React from "react";

import AppBar from '@material-ui/core/AppBar';
import ToolBar from '@material-ui/core/Toolbar';
import MenuIcon from '@material-ui/icons/Menu';

import { createStyles, withStyles, WithStyles } from '@material-ui/core';

const styles = createStyles({
  root: {},
})

export interface IMAppBarProps extends WithStyles<typeof styles> {}

/**
 * MAppBar provides lateral navigation between views within Metal. It is the
 * primary means by which users transition between different sub-applications
 * within Metal.
 */
const MAppBar: React.SFC<IMAppBarProps> = (props: IMAppBarProps): JSX.Element => {
  return (
    <div className={props.classes.root}>
      <AppBar position='static'>
        <ToolBar>
          <MenuIcon />
        </ToolBar>
      </AppBar>
    </div>
  );
};

export default withStyles(styles)(MAppBar);
