// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ThemeDefinition} from 'vuetify';
import '../styles/themes/default/main.scss';
import {black, blue, green, grey, neutral, orange, red, white, yellow} from './Colors';

const light: ThemeDefinition = {
  dark: false,
  colors: {
    mbti: yellow[45],
    primary: blue[45],
    secondary: black,
    error: red[45],
    success: green[45],
    warning: yellow[45],
    info: blue[45],
    button: neutral[55],
    font: grey[25],
    sidebar: neutral[95],

    // project status
    projectNew: orange[100],
    projectActive: blue[100],
    projectArchived: grey[30],
    projectGroup: orange[45],
    projectDeprecated: grey[30],

    // policy rule status
    prActive: blue[45],
    prInactive: grey[30],
    prDeprecated: orange[50],

    // version status
    versionNew: green[45],
    versionUnreviewed: grey[50],
    versionRejected: red[55],
    versionFreezed: grey[30],
    versionAcceptableAfterChanges: orange[60],
    versionApproved: green[45],

    borderCard: grey[80],

    tableRowOddBackgroundColor: neutral[99],
    tableRowOddBackgroundColorHover: neutral[92],
    tableRowEvenBackgroundColor: neutral[97],
    tableRowEvenBackgroundColorHover: neutral[91],
    tableHeaderBackgroundColor: neutral[97],
    tableBorderColor: grey[75],

    // modals
    successModalPlaceHolderBorderColor: green[70],
    errorModalPlaceHolderBorderColor: red[40],
    warningModalPlaceHolderBorderColor: yellow[70],
    warningPulseColor: yellow[70],
    modalBorderColor: grey[50],

    // dialog
    dialogBackground: white,
    dialogBorder: grey[30],

    // markdown
    markdownBackground: neutral[97],

    // licences
    licenceChartIcon: yellow[40],
    licenceForbidden: red[35],
    licenceNotApproved: neutral[75],

    //RuleButton
    ruleButton: grey[50],

    //labels
    labels: grey[50],

    // background
    backgroundColor: neutral[99],

    // PolicyStatus/Components Status
    policyStatusDeniedColor: red[45],
    policyStatusUnassertedColor: red[45],
    policyStatusWarnedColor: yellow[45],

    // panel / cards
    cardBackground: white,
    cardBorder: grey[75],
    errorMessageColor: red[40],

    // Charts
    chartRed: red[35],
    chartOrange: orange[50],
    chartYellow: yellow[45],
    chartGreen: green[45],
    chartGrey: grey[50],

    chartFLRed: blue[30],
    chartFLOrange: blue[35],
    chartFLYellow: blue[75],
    chartFLGreen: green[60],
    chartFLGrey: grey[50],
    chartLabelColor: black,

    tabsBackground: white,
    tabsHeader: neutral[96],
    inactiveTabColor: neutral[50],
    backgroundActiveTab: neutral[92],

    inputColor: grey[50],

    // Background Color of AND and OR
    backgroundColorAndOr: neutral[90],

    dashboardcardBackground: grey[75],

    // Notification Bar
    notificationBarBackground: yellow[90],
    notificationBarTextColor: black,

    // Background
    background: neutral[98], // Hintergrundfarbe

    // fontColorTableSecond
    fontColorTableSecond: grey[50],

    // approval status colors
    approvalApproved: green[45],
    approvalDeclined: red[35],
    approvalPending: orange[50],
  },
};
export default light;
