<?php

$darkTheme = false;

?>
<!doctype html>
<html class="no-js" lang="en">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Factoid Wallet - <?php echo $pageTitle; ?></title>
        <link rel="stylesheet" href="css/app.css">
    </head>
    <body<?php if($darkTheme == true){ echo ' class="darkTheme"'; }  ?>>
        
        <section id="frame" class="row align-stretch">
            <section class="leftCol small-12 medium-3 columns">
                <header>
                    <img src="img/factom_stacked.svg" class="svg logo" alt="Factom">
                    <h1>Factoid Wallet</h1>
                </header>
                <nav>
                    <ul>
                        <li<?php if($activeNav == 1){ echo ' class="active"'; } ?>><a href="/" class="transactions"><i><img src="img/nav_transactions.svg" class="svg" alt="Transactions"></i>Transactions</a></li>
                        <li<?php if($activeNav == 2){ echo ' class="active"'; } ?>><a href="address-book.php" class="address-book"><i><img src="img/nav_address-book.svg" class="svg" alt="Address Book"></i>Address Book</a></li>
                        <li<?php if($activeNav == 3){ echo ' class="active"'; } ?>><a href="settings.php" class="settings"><i><img src="img/nav_settings.svg" class="svg" alt="Settings"></i>Settings</a></li>
                    </ul>
                </nav>
            </section>
            <section class="rightCol small-12 medium-expand columns">
                <section class="balances">
                    <div class="row">
                        <div class="small-6 columns">
                            <section class="balance factoids dec8">
                                <i><img src="img/balance_factoids.svg" class="svg" alt="Factoids"></i>
								<span>5,123<small>.32100001</small></span>
                            </section>
                        </div>
                        <div class="small-6 columns">
                            <section class="balance credits">
                                <i><img src="img/balance_credits.svg" class="svg" alt="Entry Credits"></i>
                                <span>1,422</span>
                            </section>
                        </div>
                    </div>
                </section>
                <main class="<?php echo $mainClass ?>">
                    <?php echo $mainContent; ?>
                </main>
            </section>
        </section>
        
        <script src="bower_components/jquery/dist/jquery.min.js"></script>
        <script src="bower_components/what-input/what-input.js"></script>
        <script src="bower_components/foundation-sites/dist/foundation.js"></script>
        <script src="js/app.js"></script>
        <script>
            $(window).load(function() {
                $(document).foundation();
            });
        </script>
    </body>
</html>
