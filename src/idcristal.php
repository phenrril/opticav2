<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.1/jquery.min.js"></script>
<?php include_once "includes/header.php"; ?>

<div class="row">
    <div class="col-lg-12">
        <div class="form-group">
            <h4 class="text-center">ID Cristales</h4><br><br>
        </div>
    </div>
</div>
            <form method="post" id="form_cristal">
                <div class="row justify-content-center">
                    <div class="col-md-4 text-center">
                        <div class="card">
                            <div class="card-header">
                                Buscar ID Venta
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="idventa" class="form-control" type="text" name="idventa" placeholder="Ingresá el Id de la venta">
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4 text-center">
                        <div class="card">
                            <div class="card-header">
                                Colocar ID Cristal
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="idcristal" class="form-control" type="text" name="idcristal" placeholder="Ingresá el Id de cristales">
                                </div>
                            </div>
                        </div>
                    </div>                  
                </div>
            </form> 
<div class="row justify-content-center">
    <input type="button" class="btn btn-primary" value="Aplicar Descuento" id="guardar_cristal" name="guardar_cristal" onclick=""></input> 
</div>
<div id="div_cristal"></div>

<?php include_once "includes/footer.php"; ?>