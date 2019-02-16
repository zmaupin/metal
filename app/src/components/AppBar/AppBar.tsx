import * as React from "react";

import { createStyles, withStyles, WithStyles } from '@material-ui/core';

const styles = createStyles({
  root: {},
})

export interface IMAppBarProps extends WithStyles<typeof styles> {}

const MAppBar: React.SFC<IMAppBarProps> = (props: IMAppBarProps): JSX.Element => {
  return (
    <div className={props.classes.root}>

    </div>
  );
};

export default withStyles(styles)(MAppBar);
