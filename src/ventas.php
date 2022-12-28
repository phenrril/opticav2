
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.1/jquery.min.js"></script>
<?php include_once "includes/header.php";
require("../conexion.php");

$id_user = $_SESSION['idUser'];
$permiso = "nueva_venta";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}
?>
<div class="row">
    <div class="col-lg-12">
        <div class="form-group">
            <h4 class="text-center">Datos del Cliente</h4>
        </div>
        <div class="card">
            <div class="card-body">
                <form method="post">
                    <div class="row">
                        <div class="col-lg-4">
                            <div>
                                <input type="hidden" id="idcliente" value="1" name="idcliente" required>
                                <label>Nombre</label>
                                <input type="text" name="nom_cliente" id="nom_cliente" class="form-control" placeholder="Ingrese nombre del cliente" required>
                            </div>
                        </div>
                        <div class="col-lg-4">
                            <div class="form-group">
                                <label>Teléfono</label>
                                <input type="number" name="tel_cliente" id="tel_cliente" class="form-control" disabled required>
                            </div>
                        </div>
                        <div class="col-lg-4">
                            <div class="form-group">
                                <label>Dirreción</label>
                                <input type="text" name="dir_cliente" id="dir_cliente" class="form-control" disabled required>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        </div>
        <div class="card">
            <div class="card-header bg-primary text-white text-center">
                Datos Venta
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-lg-6">
                        <div class="form-group">
                            <i class="fas fa-user"></i><p style="font-size: 16px; text-transform: uppercase; color: red;"><?php echo $_SESSION['nombre']; ?></p>
                            <label> <b>Graduaciones:</b></label>
                            <p style="font-size: 16px; text-transform: uppercase; color: black;">
                            <table class="table table-borderless" id="tablaGracuadiones">
                            <form id="graduaciones">
                                <div id="okGrad">
                                <tr>
                                    
                                    <td><b>Graduacion Cerca </b></td>
                                </tr>
                                
                                <tr>
                                    <td><b>Ojo D: </b><input name="ojoD1" id="ojoD1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoD2" name="ojoD2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoD3" name="ojoD3" type="text" size="4"></td>
                                </tr>
                                <tr>
                                    <td><b>Ojo I: &nbsp</b><input id="ojoI1" name="ojoI1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoI2" name="ojoI2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoI3" name="ojoI3" type="text" size="4"></td>
                                </tr>
                                <tr>
                                    <td><b>Graduacion Lejos </b></td>
                                </tr>
                                <tr>
                                    <td><div class="col-xs-3"><b>Ojo D: </b><input  id="ojoDl1" name="ojoDl1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoDl2" name="ojoDl2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoDl3" name="ojoDl3" type="text" size="4"></div></td>
                                </tr>
                                <tr>
                                    <td><b>Ojo I: &nbsp</b><input id="ojoIl1" name="ojoIl1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoIl2" name="ojoIl2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoIl3" name="ojoIl3" type="text" size="4"></td>
                                </tr>
                                <tr>
                                    <td><input class="btn btn-primary" name="grad" id="grad" type="button" value="Agregar"></td>
                                &nbsp&nbsp&nbsp<b>ADD: &nbsp</b><input id="add" name="add" type="text" size="4">&nbsp&nbsp&nbsp<b>Obs: &nbsp</b><input id="add" name="add" type="text" size="11"></td>
                                </tr>
                                </div>
                            </form>
                        </table>    
                        </p>
                        </div>
                    </div>
                    <div class="col-lg-6">
                        <div class="card">
                            <div class="card-header">
                                Buscar Producto
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="producto" class="form-control" type="text" name="producto" placeholder="Ingresa el código o nombre">
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
        <div class="table-responsive">
            <table class="table table-hover" id="tblDetalle">
                <thead class="thead-dark">
                    <tr>
                        <th>Id</th>
                        <th>Descripción</th>
                        <th>Cantidad</th>
                        <th>Precio</th>
                        <th>Precio Total</th>
                        <th>Accion</th>
                    </tr>
                </thead>
                <tbody id="detalle_venta">

                </tbody>
                <tfoot>
                    <tr class="font-weight-bold">
                        <td colspan=3>Total Pagar</td>
                        <td></td>
                    </tr>
                </tfoot>
            </table>

        </div>
    </div>
    <div class="col-md-6">
        <a href="#" class="btn btn-primary" id="btn_generar"><i class="fas fa-save"></i> Generar Venta</a>
    </div>

</div>
<script >
    
    // $("#grad").click(function(){
    //         $.ajax({
    //                 url: "resultado.php",
    //                 type: "POST",
    //                 data: $("#graduaciones").serialize(),
    //                 success: function(resultado){
    //                         $("#okGrad").html(resultado);
    //                 }
    //         });
    //});
    
</script>
<?php include_once "includes/footer.php"; ?>